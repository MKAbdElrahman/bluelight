package handlers

import (
	"log/slog"
	"net/http"
)

type RouterConfig struct {
	API_Environment string
	API_Version     string
	Logger          *slog.Logger
}

func NewRouter(cfg RouterConfig) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /v1/healthcheck", newHealthCheckHandlerFunc(cfg.API_Environment, cfg.API_Version))
	router.HandleFunc("POST /v1/movies", newCreateMovieHandlerFunc())
	router.HandleFunc("GET /v1/movies/{id}", newShowMovieHandlerFunc())

	return router
}
