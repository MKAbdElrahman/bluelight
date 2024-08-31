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

func NewCreateAuthTokenHandlerFunc(em *errorhandler.ErrorHandeler, userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewCreateAuthTokenRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}

		// Business
		params := user.CreateAuthTokenParams{
			Email:    req.Body.Email,
			Password: req.Body.Password,
		}

		t, err := userService.CreateAuthToken(params)
		if err != nil {

			var validErr *verrors.ValidationError
			switch {

			case errors.As(err, &validErr):
				em.SendClientError(w, r, apierror.UnprocessableEntityError.WithValidationError(validErr))
			case errors.Is(err, user.ErrRecordNotFound):
				em.SendClientError(w, r, apierror.UnauthorizedError)
			case errors.Is(err, user.ErrInvalidCredentials):
				em.SendClientError(w, r, apierror.UnauthorizedError)
			default:
				em.SendServerError(w, r, apierror.NewInternalServerError(err))
			}
			return
		}
		res := v1.CreateAuthResponse{
			Expiry: t.Expiry,
			Token:  t.Plaintext,
		}

		err = jsonio.SendJSON(w, jsonio.Envelope{"token": res}, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, apierror.NewInternalServerError(err))

			return
		}
	}
}
