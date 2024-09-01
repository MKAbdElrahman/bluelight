package middleware

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	userhandlers "bluelight.mkcodedev.com/src/api/handlers/user"
)

func RequireAuthenticatedUser(em *errorhandler.ErrorHandeler) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				u := userhandlers.GetUserFromContext(r)
				if u.IsAnonymous() {
					em.SendClientError(w, r, &apierror.ClientError{
						Code:              http.StatusUnauthorized,
						UserFacingMessage: "you must be authenticated to access this resource",
					})
					return
				}
				next.ServeHTTP(w, r)
			})
	}
}
