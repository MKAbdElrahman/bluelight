package errorhandler

import (
	"errors"
	"log/slog"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
)

// ErrorManager handles error logging and responses.
type ErrorHandeler struct {
	Logger          *slog.Logger
	LogClientErrors bool
	LogServerErrors bool
}

func NewErrorHandler(logger *slog.Logger) *ErrorHandeler {
	return &ErrorHandeler{
		Logger:          logger,
		LogClientErrors: false,
		LogServerErrors: true,
	}
}

// SendServerError handles server-side errors.
func (e *ErrorHandeler) SendServerError(w http.ResponseWriter, r *http.Request, serverErr *apierror.ServerError) {
	if e.LogServerErrors {
		logError(e.Logger, r, errors.New(serverErr.InternalMessage), serverErr.Code)
	}

	data := jsonio.Envelope{"error": http.StatusText(serverErr.Code)}
	err := jsonio.SendJSON(w, data, serverErr.Code, nil)

	if err != nil {
		if e.LogServerErrors {
			logError(e.Logger, r, err, http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// SendClientError handles client-side errors.
func (e *ErrorHandeler) SendClientError(w http.ResponseWriter, r *http.Request, clientErr *apierror.ClientError) {
	if e.LogClientErrors {
		logError(e.Logger, r, errors.New(clientErr.UserFacingMessage), clientErr.Code)
	}

	data := jsonio.Envelope{"error": clientErr}
	err := jsonio.SendJSON(w, data, clientErr.Code, nil)
	if err != nil {
		if e.LogServerErrors {
			logError(e.Logger, r, err, http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func logError(logger *slog.Logger, r *http.Request, err error, statusCode int) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	if statusCode < 500 {
		logger.Error(err.Error(), "type", "client error", "method", method, "uri", uri, "status_code", statusCode)
		return
	}
	logger.Error(err.Error(), "type", "server error", "method", method, "uri", uri, "status_code", statusCode)

}
