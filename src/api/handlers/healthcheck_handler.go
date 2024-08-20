package handlers

import (
	"log/slog"
	"net/http"
)

func newHealthCheckHandlerFunc(logger *slog.Logger, env, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		health := struct {
			Status      string `json:"status"`
			Environment string `json:"environment"`
			Version     string `json:"version"`
		}{
			Status:      "available",
			Environment: env,
			Version:     version,
		}
		err := sendJSON(w, envelope{"health_check": health}, http.StatusOK, nil)
		if err != nil {
			internal(logger, w, r, err)
			return
		}
	}
}
