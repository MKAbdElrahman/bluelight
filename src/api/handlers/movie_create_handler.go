package handlers

import (
	"fmt"
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
			badRequest(logger, w, r, err.Error())
		}
		fmt.Println(request.Body)
	}
}
