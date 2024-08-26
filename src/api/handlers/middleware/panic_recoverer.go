package middleware

import (
	"fmt"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	 "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
)

func PanicRecoverer(em *errorhandler.ErrorHandeler) middlewareFunc {
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
