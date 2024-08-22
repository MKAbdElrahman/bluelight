package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RouterConfig struct {
	API_Environment string
	API_Version     string
	Logger          *slog.Logger
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	// MIDDLEWARE
	r.Use(requestSizeLimiter(1_048_576)) // 1MB
	r.Use(requestLogger(cfg.Logger))
	r.Use(panicRecoverer(cfg.Logger))

	// INVALID ROUTES HANDLERS
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		sendNotFoundError(cfg.Logger, w, r)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		sendMethodNotAllowedError(cfg.Logger, w, r)
	})
	// ROUTES
	r.Get("/v1/healthcheck", newHealthCheckHandlerFunc(cfg.Logger, cfg.API_Environment, cfg.API_Version))
	r.Post("/v1/movies", newCreateMovieHandlerFunc(cfg.Logger))
	r.Get("/v1/movies/{id}", newShowMovieHandlerFunc(cfg.Logger))

	return r
}
