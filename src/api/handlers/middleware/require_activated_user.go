package middleware

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	userhandlers "bluelight.mkcodedev.com/src/api/handlers/user"
)

func RequireActivatedUser(em *errorhandler.ErrorHandeler) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return RequireAuthenticatedUser(em)(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				u := userhandlers.GetUserFromContext(r)
				if !u.Activated {
					em.SendClientError(w, r, &apierror.ClientError{
						Code:              http.StatusForbidden,
						UserFacingMessage: "your user account must be activated to access this resource",
					})
					return
				}
				next.ServeHTTP(w, r)
			}))
	}
}
