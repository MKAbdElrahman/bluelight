package handlers

import (
	"log/slog"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newCreateMovieHandlerFunc(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request v1.CreateMovieRequest
		err := jsonio.NewJSONReader().ReadJSON(r, &request.Body)
		if err != nil {
			sendClientError(logger, w, r, v1.BadRequestError.WithDetails(map[string]string{
				"error": err.Error(),
			}))
			return
		}
		vErrors := request.Validate()
		if vErrors.Length() > 0 {
			sendClientError(logger, w, r, v1.UnprocessableEntityError.WithDetails(vErrors.Details()))
			return
		}
	}
}
