package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

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

func RequestLogger(logger *slog.Logger) middlewareFunc {
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
			logger.Info("HTTP Request Handled",
				append(logFields, slog.Int("code", crw.statusCode), slog.String("status", http.StatusText(crw.statusCode)))...)

		})
	}
}
