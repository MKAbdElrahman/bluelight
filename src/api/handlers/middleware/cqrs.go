package middleware

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
)

func CQRS(em *errorhandler.ErrorHandeler, trustedOrigins []string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				w.Header().Add("Vary", "Origin")
				w.Header().Add("Vary", "Access-Control-Request-Method")

				origin := r.Header.Get("Origin")

				if origin != "" {
					for i := range trustedOrigins {
						if origin == trustedOrigins[i] {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {

								w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
								w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

								w.WriteHeader(http.StatusOK)
								return
							}

							break
						}
					}

				}
				next.ServeHTTP(w, r)
			})
	}
}
