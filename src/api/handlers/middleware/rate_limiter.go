package middleware

import (
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"golang.org/x/time/rate"
)

func RateLimiter(em *errorhandler.ErrorHandeler) middlewareFunc {
	return func(next http.Handler) http.Handler {
		limiter := rate.NewLimiter(2, 4)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				em.SendClientError(w, r, v1.TooManyRequestsError)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
