package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/handlers/middleware"
	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/infrastructure/db/repositories"
	"github.com/go-chi/chi/v5"
)

type RouterConfig struct {
	API_Environment string
	API_Version     string
	Logger          *slog.Logger
	DB              *sql.DB
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	em := errorhandler.NewErrorHandler(cfg.Logger)
	em.LogClientErrors = true
	em.LogServerErrors = true

	// MIDDLEWARE
	r.Use(middleware.RequestSizeLimiter(1_048_576)) // 1MB
	r.Use(middleware.RequestLogger(cfg.Logger))
	r.Use(middleware.PanicRecoverer(em))

	// INVALID ROUTES HANDLERS
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		em.SendClientError(w, r, v1.NotFoundError)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		em.SendClientError(w, r, v1.MethodNotAllowedError)
	})

	// INFRASTRUCTURE

	movieRepository := repositories.NewPostgresMovieRepository(cfg.DB)
	movieService := domain.NewMovieService(movieRepository)

	// ROUTES
	r.Get("/v1/healthcheck", newHealthCheckHandlerFunc(em, cfg.API_Environment, cfg.API_Version))
	r.Post("/v1/movies", newCreateMovieHandlerFunc(em, movieService))
	r.Patch("/v1/movies/{id}", newUpdateMovieHandlerFunc(em, movieService))
	r.Get("/v1/movies/{id}", newShowMovieHandlerFunc(em, movieService))
	r.Delete("/v1/movies/{id}", newDeleteMovieHandlerFunc(em, movieService))

	return r
}
