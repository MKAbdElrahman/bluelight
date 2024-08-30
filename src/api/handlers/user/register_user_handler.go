package userhandlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	v1 "bluelight.mkcodedev.com/src/api/contracts/v1/user"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
	"bluelight.mkcodedev.com/src/core/domain/user"
	"bluelight.mkcodedev.com/src/core/domain/verrors"
	"bluelight.mkcodedev.com/src/infrastructure/mailer"
)

func NewRegisterUserHandlerFunc(backgroundRoutinesWaitGroup *sync.WaitGroup, em *errorhandler.ErrorHandeler, userService *user.UserService, mailerService *mailer.Mailer) http.HandlerFunc {
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

		backgroundRoutinesWaitGroup.Add(1)
		background(em.Logger, func() {
			defer backgroundRoutinesWaitGroup.Done()
			err = mailerService.WelcomeNewRegisteredUser(context.Background(), u.Email, u.Name)
			if err != nil {
				em.Logger.Error("failed to send welcome email after retries", "err", err)
			}
		})

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

func background(logger *slog.Logger, fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic", "err", fmt.Sprintf("%v", err))
			}
		}()
		fn()
	}()

}
