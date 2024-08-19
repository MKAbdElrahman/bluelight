package handlers

import (
	"log/slog"
	"net/http"
)

type MuxConfig struct {
	Environment string
	Version     string
	Logger      *slog.Logger
}

func NewRouter(cfg MuxConfig) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /v1/healthcheck", newHealthCheckHandlerFunc(cfg.Environment, cfg.Version))
	router.HandleFunc("POST /v1/movies", newCreateMovieHandlerFunc())
	router.HandleFunc("GET /v1/movies/{id}", newShowMovieHandlerFunc())

	return router
}
