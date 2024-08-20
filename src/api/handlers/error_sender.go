package handlers

import (
	"log/slog"
	"net/http"
)

// internal responds with a 500 Internal Server Error error.
func internal(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	logError(logger, r, err)
	message := http.StatusText(http.StatusInternalServerError)
	sendJSONError(logger, w, r, http.StatusInternalServerError, message)
}

// notFound responds with a 404 Not Found error.
func notFound(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	message := http.StatusText(http.StatusNotFound)
	sendJSONError(logger, w, r, http.StatusNotFound, message)
}

// methodNotAllowed responds with a 405 Method Not Allowed error.
func methodNotAllowed(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	message := http.StatusText(http.StatusMethodNotAllowed)
	sendJSONError(logger, w, r, http.StatusMethodNotAllowed, message)
}

func sendJSONError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, status int, message any) {
	data := envelope{"error": message}
	err := sendJSON(w, data, status, nil)
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
