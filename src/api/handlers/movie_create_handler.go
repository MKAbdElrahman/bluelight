package handlers

import (
	"log/slog"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/jsonio"
)

func newCreateMovieHandlerFunc(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request v1.CreateMovieRequest
		err := jsonio.NewJSONReader().ReadJSON(r, &request.Body)
		if err != nil {
			sendBadRequestError(logger, w, r, err)
			return
		}
		vErrors := request.Validate()
		if vErrors.Length() > 0 {
			sendValidationError(logger, w, r, vErrors)
			return
		}
	}
}
