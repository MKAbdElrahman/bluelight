package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/errormanager"
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

	em := errormanager.NewErrorManager(cfg.Logger)
	em.LogClientErrors = true
	em.LogServerErrors = true

	// MIDDLEWARE
	r.Use(requestSizeLimiter(1_048_576)) // 1MB
	r.Use(requestLogger(cfg.Logger))
	r.Use(panicRecoverer(em))

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
	r.Get("/v1/movies/{id}", newShowMovieHandlerFunc(em))

	return r
}
