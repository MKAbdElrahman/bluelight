package errorhandler

import (
	"fmt"
	"log/slog"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

// ErrorManager handles error logging and responses.
type ErrorHandeler struct {
	logger          *slog.Logger
	LogClientErrors bool
	LogServerErrors bool
}

func NewErrorHandler(logger *slog.Logger) *ErrorHandeler {
	return &ErrorHandeler{
		logger:          logger,
		LogClientErrors: false,
		LogServerErrors: true,
	}
}

// SendServerError handles server-side errors.
func (e *ErrorHandeler) SendServerError(w http.ResponseWriter, r *http.Request, serverErr v1.ServerError) {
	if e.LogServerErrors {
		logError(e.logger, r, fmt.Errorf(serverErr.InternalMessage), serverErr.Code)
	}

	data := jsonio.Envelope{"error": http.StatusText(serverErr.Code)}
	err := jsonio.SendJSON(w, data, serverErr.Code, nil)

	if err != nil {
		if e.LogServerErrors {
			logError(e.logger, r, err, http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// SendClientError handles client-side errors.
func (e *ErrorHandeler) SendClientError(w http.ResponseWriter, r *http.Request, clientErr *v1.ClientError) {
	if e.LogClientErrors {
		logError(e.logger, r, fmt.Errorf(clientErr.UserFacingMessage), clientErr.Code)
	}

	data := jsonio.Envelope{"error": clientErr}
	err := jsonio.SendJSON(w, data, clientErr.Code, nil)
	if err != nil {
		if e.LogServerErrors {
			logError(e.logger, r, err, http.StatusInternalServerError)
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
