package handlers

import (
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newHealthCheckHandlerFunc(eh *errorhandler.ErrorHandeler, env, version string) http.HandlerFunc {
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
			eh.SendServerError(w, r, v1.InternalServerError)
			return
		}
	}
}
