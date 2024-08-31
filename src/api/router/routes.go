package router

import (
	"database/sql"
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
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	em := errorhandler.NewErrorHandler(cfg.Logger)
	em.LogClientErrors = true
	em.LogServerErrors = true

	// MIDDLEWARE
	r.Use(middleware.RateLimiter(em, cfg.LimiterConfig))
	r.Use(middleware.RequestSizeLimiter(1_048_576)) // 1MB
	r.Use(middleware.RequestLogger(cfg.Logger))
	r.Use(middleware.PanicRecoverer(em))

	// INVALID ROUTES HANDLERS
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		em.SendClientError(w, r, apierror.NotFoundError)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		em.SendClientError(w, r, apierror.MethodNotAllowedError)
	})

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

	userService := user.NewUserService(userRepository, tokenRepository, cfg.Mailer)
	// ROUTES
	r.Get("/v1/healthcheck", healthCheckHandlers.NewHealthCheckHandlerFunc(em, cfg.API_Environment, cfg.API_Version))
	r.Post("/v1/movies", movieHandlers.NewCreateMovieHandlerFunc(em, movieService))
	r.Patch("/v1/movies/{id}", movieHandlers.NewUpdateMovieHandlerFunc(em, movieService))
	r.Get("/v1/movies/{id}", movieHandlers.NewShowMovieHandlerFunc(em, movieService))
	r.Get("/v1/movies", movieHandlers.NewListMovieHandlerFunc(em, movieService))
	r.Delete("/v1/movies/{id}", movieHandlers.NewDeleteMovieHandlerFunc(em, movieService))

	r.Post("/v1/users", userHandlers.NewRegisterUserHandlerFunc(cfg.BackgroundWaitGroup, em, userService))
	r.Put("/v1/users/activate", userHandlers.NewActivateUserHandlerFunc(cfg.BackgroundWaitGroup, em, userService))
	r.Post("/v1/tokens/authentication", userHandlers.NewCreateAuthTokenHandlerFunc(em, userService))

	return r
}
