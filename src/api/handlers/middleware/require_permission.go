package middleware

import (
	"fmt"
	"net/http"

	"bluelight.mkcodedev.com/src/core/domain/user"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	userhandlers "bluelight.mkcodedev.com/src/api/handlers/user"
)

func RequirePermission(em *errorhandler.ErrorHandeler, userService *user.UserService, code string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return RequireActivatedUser(em)(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				u := userhandlers.GetUserFromContext(r)

				permissions, err := userService.GetAllPermissionsForUser(u.Id)
				if err != nil {
					em.SendServerError(w, r, &apierror.InternalServerError)
					return
				}
				fmt.Println(permissions)
				if !permissions.Include(code) {
					em.SendClientError(w, r, &apierror.ClientError{
						Code:              http.StatusForbidden,
						UserFacingMessage: "your user account doesn't have the necessary permissions to access this resource",
					})
					return
				}
				next.ServeHTTP(w, r)
			}))
	}
}
