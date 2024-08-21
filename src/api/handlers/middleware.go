package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type middleware func(next http.Handler) http.Handler

func panicRecoverer(logger *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						w.Header().Set("Connection", "close")
						internal(logger, w, r, fmt.Errorf("%s", err))
					}
				}()
				next.ServeHTTP(w, r)
			})
	}
}

func requestLogger(logger *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Info("HTTP Request Handled",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.Duration("duration", time.Since(start)),
			)
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
