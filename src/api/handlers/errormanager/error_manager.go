package errormanager

import (
	"fmt"
	"log/slog"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

// ErrorManager handles error logging and responses.
type ErrorManager struct {
	logger          *slog.Logger
	LogClientErrors bool
	LogServerErrors bool
}

// NewErrorManager creates a new ErrorManager with customizable logging options.
func NewErrorManager(logger *slog.Logger) *ErrorManager {
	return &ErrorManager{
		logger:          logger,
		LogClientErrors: false,
		LogServerErrors: true,
	}
}

// SendServerError handles server-side errors.
func (e *ErrorManager) SendServerError(w http.ResponseWriter, r *http.Request, serverErr v1.ServerError) {
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
func (e *ErrorManager) SendClientError(w http.ResponseWriter, r *http.Request, clientErr v1.ClientError) {
	if e.LogClientErrors {
		logWarn(e.logger, r, fmt.Errorf(clientErr.UserFacingMessage), clientErr.Code)
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
	logger.Error(err.Error(), "method", method, "uri", uri, "status_code", statusCode)
}

func logWarn(logger *slog.Logger, r *http.Request, err error, statusCode int) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	logger.Warn(err.Error(), "method", method, "uri", uri, "status_code", statusCode)
}
