package userhandlers

import (
	"errors"
	"net/http"
	"sync"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	v1 "bluelight.mkcodedev.com/src/api/contracts/v1/user"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
	"bluelight.mkcodedev.com/src/core/domain/user"
	"bluelight.mkcodedev.com/src/core/domain/verrors"
)

func NewActivateUserHandlerFunc(backgroundRoutinesWaitGroup *sync.WaitGroup, em *errorhandler.ErrorHandeler, userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewActivateUserRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}

		// Business
		params := user.UserActivationParams{
			TokenPlaintext: req.Body.TokenPlaintext,
		}

		u, err := userService.ActivateUser(backgroundRoutinesWaitGroup, em.Logger, params)
		if err != nil {

			var validErr *verrors.ValidationError
			switch {

			case errors.As(err, &validErr):
				em.SendClientError(w, r, apierror.UnprocessableEntityError.WithValidationError(validErr))
			case errors.Is(err, user.ErrEditConflict):
				em.SendClientError(w, r, apierror.ConflictError)
			case errors.Is(err, user.ErrRecordNotFound):
				em.SendClientError(w, r, apierror.NotFoundError)
			default:
				em.SendServerError(w, r, apierror.NewInternalServerError(err))
			}
			return
		}
		res := v1.ActivateUserResponse{
			Activated: u.Activated,
		}

		err = jsonio.SendJSON(w, jsonio.Envelope{"user": res}, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, apierror.NewInternalServerError(err))

			return
		}
	}
}
