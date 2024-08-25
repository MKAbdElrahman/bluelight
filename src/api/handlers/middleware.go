package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/errormanager"
)

type middleware func(next http.Handler) http.Handler

func panicRecoverer(em *errormanager.ErrorManager) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						w.Header().Set("Connection", "close")
						em.SendServerError(w, r, v1.ServerError{
							InternalMessage: fmt.Sprintf("%s", err),
							Code:            http.StatusInternalServerError,
						})
					}
				}()
				next.ServeHTTP(w, r)
			})
	}
}

// RequestSize is a middleware that will limit request sizes to a specified
// number of bytes. It uses MaxBytesReader to do so.
func requestSizeLimiter(bytes int64) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, bytes)
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}

// customResponseWriter is a wrapper around http.ResponseWriter to capture the status code.
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code.
func (w *customResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func requestLogger(logger *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			crw := &customResponseWriter{w, http.StatusOK}

			next.ServeHTTP(crw, r)

			duration := time.Since(start)
			logFields := []any{
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.Duration("duration", duration),
			}

			if crw.statusCode >= 500 {
				logger.Error("HTTP Request Handled with Server Error",
					append(logFields, slog.Int("status", crw.statusCode))...)
			} else if crw.statusCode >= 400 {
				logger.Warn("HTTP Request Handled with Client Error",
					append(logFields, slog.Int("status", crw.statusCode))...)
			} else {
				logger.Info("HTTP Request Handled Successfully",
					append(logFields, slog.Int("status", crw.statusCode))...)
			}
		})
	}
}
