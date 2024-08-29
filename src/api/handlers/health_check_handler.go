package handlers

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
)

func newHealthCheckHandlerFunc(em *errorhandler.ErrorHandeler, env, version string) http.HandlerFunc {
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
			em.SendServerError(w, r, apierror.NewInternalServerError(err))
			return
		}
	}
}
