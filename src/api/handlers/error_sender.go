package handlers

import (
	"log/slog"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/jsonio"
)

func sendAPIError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, apiErr v1.ApiError) {
	data := jsonio.Envelope{"error": apiErr}
	err := jsonio.SendJSON(w, data, apiErr.Code, nil)
	if err != nil {
		logError(logger, r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func sendInternalError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	logError(logger, r, err)
	message := http.StatusText(http.StatusInternalServerError)
	sendAPIError(logger, w, r, v1.ApiError{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

func sendNotFoundError(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	sendAPIError(logger, w, r, v1.ApiError{
		Code:    code,
		Message: http.StatusText(code),
	})
}

func sendMethodNotAllowedError(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	code := http.StatusMethodNotAllowed
	sendAPIError(logger, w, r, v1.ApiError{
		Code:    code,
		Message: http.StatusText(code),
	})
}

func sendBadRequestError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	code := http.StatusBadRequest
	sendAPIError(logger, w, r, v1.ApiError{
		Code:    code,
		Message: err.Error(),
	})
}

func sendValidationError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, errors v1.ValidationError) {
	code := http.StatusUnprocessableEntity
	sendAPIError(logger, w, r, v1.ApiError{
		Code:    code,
		Message: http.StatusText(code),
		Details: errors.Errors,
	})
}

func logError(logger *slog.Logger, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	logger.Error(err.Error(), "method", method, "uri", uri)
}
