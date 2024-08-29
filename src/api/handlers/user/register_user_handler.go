package userhandlers

import (
	"errors"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	v1 "bluelight.mkcodedev.com/src/api/contracts/v1/user"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
	"bluelight.mkcodedev.com/src/core/domain/user"
	"bluelight.mkcodedev.com/src/core/domain/verrors"
)

func NewRegisterUserHandlerFunc(em *errorhandler.ErrorHandeler, userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewRegisterUserRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}

		// Business
		params := user.UserRegisterationParams{
			Name:     req.Body.Name,
			Email:    req.Body.Email,
			Password: req.Body.Password,
		}

		u, err := userService.RegisterUser(params)
		if err != nil {

			var validErr *verrors.ValidationError
			switch {
			case errors.As(err, &validErr):
				em.SendClientError(w, r, apierror.UnprocessableEntityError.WithValidationError(validErr))
			case errors.Is(err, user.ErrDuplicateEmail):
				em.SendClientError(w, r, apierror.UnprocessableEntityError.WithError("email", user.ErrDuplicateEmail))
			default:
				em.SendServerError(w, r, apierror.NewInternalServerError(err))
			}
			return
		}

		res := v1.RegisterUserResponse{
			Id:        u.Id,
			Name:      u.Name,
			Email:     u.Email,
			Activated: u.Activated,
			Version:   u.Version,
		}

		err = jsonio.SendJSON(w, jsonio.Envelope{"user": res}, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, apierror.NewInternalServerError(err))

			return
		}
	}
}
