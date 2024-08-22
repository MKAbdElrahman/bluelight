package handlers

import (
	"log/slog"
	"net/http"

	"bluelight.mkcodedev.com/src/api/handlers/jsonio"
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
		err := jsonio.SendJSON(w, jsonio.Envelope{"health_check": health}, http.StatusOK, nil)
		if err != nil {
			sendInternalError(logger, w, r, err)
			return
		}
	}
}
