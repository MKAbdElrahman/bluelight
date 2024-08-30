package user

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Mailer interface {
	WelcomeNewRegisteredUser(ctx context.Context, recipientEmail, recipientName string) error
}

type UserService struct {
	userRepository UserRepositoty
	mailerService  Mailer
}

func NewUserService(r UserRepositoty, mailerService Mailer) *UserService {
	return &UserService{
		userRepository: r,
		mailerService:  mailerService,
	}
}

type UserRegisterationParams struct {
	Name     string
	Email    string
	Password string
}

func (svc *UserService) RegisterUser(backgroundRoutinesWaitGroup *sync.WaitGroup, logger *slog.Logger, params UserRegisterationParams) (*User, error) {
	u, err := NewUser(params.Name, params.Email, params.Password)
	if err != nil {
		return nil, err
	}

	err = svc.userRepository.Create(u)
	if err != nil {
		return nil, err
	}
	backgroundRoutinesWaitGroup.Add(1)
	background(logger, func() {
		defer backgroundRoutinesWaitGroup.Done()
		err = svc.mailerService.WelcomeNewRegisteredUser(context.Background(), u.Email, u.Name)
		if err != nil {
			logger.Error("failed to send welcome email after retries", "err", err)
		}
	})
	return u, nil
}

func (svc *UserService) UpdateUser(u *User) error {
	return svc.userRepository.Update(u)
}

func (svc *UserService) GetByEmail(email string) (*User, error) {
	return svc.userRepository.GetByEmail(email)
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
