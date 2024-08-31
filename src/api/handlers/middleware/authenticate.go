package middleware

import (
	"errors"
	"net/http"
	"strings"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	userhandlers "bluelight.mkcodedev.com/src/api/handlers/user"
	"bluelight.mkcodedev.com/src/core/domain/user"
)

func Authenticate(em *errorhandler.ErrorHandeler, userService *user.UserService) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Vary", "Authorization")
				authorizationHeader := r.Header.Get("Authorization")
				if authorizationHeader == "" {
					r = userhandlers.StoreUserInContext(user.AnonymousUser, r)
					next.ServeHTTP(w, r)
					return
				}
				headerParts := strings.Split(authorizationHeader, " ")
				if len(headerParts) != 2 || headerParts[0] != "Bearer" {
					w.Header().Set("WWW-Authenticate", "Bearer")
					em.SendClientError(w, r, &apierror.ClientError{
						Code:              http.StatusUnauthorized,
						UserFacingMessage: "invalid or missing authentication token",
					})
					return
				}
				token := headerParts[1]
				if err := user.ValidatePlainTextForm(token); err != nil {
					em.SendClientError(w, r, &apierror.ClientError{
						Code:              http.StatusUnauthorized,
						UserFacingMessage: "invalid or missing authentication token",
					})
					return
				}

				u, err := userService.GetUserByToken(user.ScopeAuthentication, token)
				if err != nil {
					switch {
					case errors.Is(err, user.ErrRecordNotFound):
						em.SendClientError(w, r, &apierror.ClientError{
							Code:              http.StatusUnauthorized,
							UserFacingMessage: "invalid or missing authentication token",
						})
					default:
						em.SendServerError(w, r, &apierror.InternalServerError)

					}
					return
				}

				r = userhandlers.StoreUserInContext(u, r)
				next.ServeHTTP(w, r)
			})
	}
}
