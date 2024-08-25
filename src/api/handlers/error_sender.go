package handlers

import (
	"log/slog"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

// sendServerError handles server-side errors.
func sendServerError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, serverErr v1.ServerError) {
	data := jsonio.Envelope{"error": serverErr}
	err := jsonio.SendJSON(w, data, serverErr.Code, nil)
	if err != nil {
		logError(logger, r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// sendClientError handles client-side errors.
func sendClientError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, clientErr v1.ClientError) {
	data := jsonio.Envelope{"error": clientErr}
	err := jsonio.SendJSON(w, data, clientErr.Code, nil)
	if err != nil {
		logError(logger, r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func logError(logger *slog.Logger, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	logger.Error(err.Error(), "method", method, "uri", uri)
}
