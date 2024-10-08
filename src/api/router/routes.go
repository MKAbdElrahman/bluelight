package router

import (
	"database/sql"
	"expvar"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	healthCheckHandlers "bluelight.mkcodedev.com/src/api/handlers/healthcheck"
	movieHandlers "bluelight.mkcodedev.com/src/api/handlers/movie"
	userHandlers "bluelight.mkcodedev.com/src/api/handlers/user"

	"bluelight.mkcodedev.com/src/api/handlers/middleware"
	"bluelight.mkcodedev.com/src/core/domain/movie"
	"bluelight.mkcodedev.com/src/core/domain/user"
	"bluelight.mkcodedev.com/src/infrastructure/db/repositories"
	"bluelight.mkcodedev.com/src/infrastructure/mailer"
	"github.com/go-chi/chi/v5"
)

type RouterConfig struct {
	API_Environment     string
	API_Version         string
	Logger              *slog.Logger
	DB                  *sql.DB
	Mailer              *mailer.Mailer
	LimiterConfig       middleware.RateLimiterConfig
	BackgroundWaitGroup *sync.WaitGroup
	TrustedOrigins      []string
}

func NewRouter(cfg RouterConfig) http.Handler {

	// INFRASTRUCTURE

	movieRepository := repositories.NewPostgresMovieRepository(
		cfg.DB,
		repositories.PostgresMovieRepositryConfig{
			Timeout: 3 * time.Second,
		})
	movieService := movie.NewMovieService(movieRepository)

	tokenRepository := repositories.NewPostgresTokenRepository(
		cfg.DB,
		repositories.PostgresTokenRepositryConfig{
			Timeout: 3 * time.Second,
		})

	userRepository := repositories.NewPostgresUserRepository(
		cfg.DB,
		repositories.PostgresUserRepositryConfig{
			Timeout: 3 * time.Second,
		})

	permissionsRepository := repositories.NewPostgresPermissionRepository(cfg.DB,
		repositories.PostgresPermissionRepositryConfig{
			Timeout: 3 * time.Second,
		})

	userService := user.NewUserService(userRepository, permissionsRepository, tokenRepository, cfg.Mailer)

	/////////

	r := chi.NewRouter()

	em := errorhandler.NewErrorHandler(cfg.Logger)
	em.LogClientErrors = true
	em.LogServerErrors = true

	// MIDDLEWARE
	r.Use(middleware.PanicRecoverer(em))
	r.Use(middleware.Metrics(em))
	r.Use(middleware.RequestLogger(cfg.Logger))
	r.Use(middleware.CQRS(em, cfg.TrustedOrigins))
	r.Use(middleware.RateLimiter(em, cfg.LimiterConfig))
	r.Use(middleware.RequestSizeLimiter(1_048_576)) // 1MB
	r.Use(middleware.Authenticate(em, userService))

	// INVALID ROUTES HANDLERS
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		em.SendClientError(w, r, apierror.NotFoundError)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		em.SendClientError(w, r, apierror.MethodNotAllowedError)
	})

	// ROUTE

	// Health check route
	r.Get("/v1/healthcheck", healthCheckHandlers.NewHealthCheckHandlerFunc(em, cfg.API_Environment, cfg.API_Version))

	// EXPVAR route
	r.Handle("/debug/vars", expvar.Handler())

	// MOVIE ROUTES
	r.Route("/v1/movies", func(r chi.Router) {
		r.Use(middleware.RequireActivatedUser(em))

		// Write endpoints (require write permission)
		r.With(middleware.RequirePermission(em, userService, "movies:write")).Post("/", movieHandlers.NewCreateMovieHandlerFunc(em, movieService))
		r.With(middleware.RequirePermission(em, userService, "movies:write")).Patch("/{id}", movieHandlers.NewUpdateMovieHandlerFunc(em, movieService))
		r.With(middleware.RequirePermission(em, userService, "movies:write")).Delete("/{id}", movieHandlers.NewDeleteMovieHandlerFunc(em, movieService))

		// Read endpoints (require read permission)
		r.With(middleware.RequirePermission(em, userService, "movies:read")).Get("/{id}", movieHandlers.NewShowMovieHandlerFunc(em, movieService))
		r.With(middleware.RequirePermission(em, userService, "movies:read")).Get("/", movieHandlers.NewListMovieHandlerFunc(em, movieService))

	})

	// User routes
	r.Post("/v1/users", userHandlers.NewRegisterUserHandlerFunc(cfg.BackgroundWaitGroup, em, userService))
	r.Put("/v1/users/activate", userHandlers.NewActivateUserHandlerFunc(cfg.BackgroundWaitGroup, em, userService))
	r.Post("/v1/tokens/authentication", userHandlers.NewCreateAuthTokenHandlerFunc(em, userService))

	return r
}
