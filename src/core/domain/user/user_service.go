package user

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"bluelight.mkcodedev.com/src/core/domain/token"
)

type Mailer interface {
	WelcomeNewRegisteredUser(ctx context.Context, u *User, activationToken string) error
}

type TokenService interface {
	New(userID int64, ttl time.Duration, scope string) (*token.Token, error)
}

type UserRepositoty interface {
	Create(u *User) error
	GetByEmail(email string) (*User, error)
	Update(u *User) error
}

type UserService struct {
	userRepository UserRepositoty
	mailerService  Mailer
	tokenService   TokenService
}

func NewUserService(r UserRepositoty, tokenService TokenService, mailerService Mailer) *UserService {

	return &UserService{
		userRepository: r,
		mailerService:  mailerService,
		tokenService:   tokenService,
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

	token, err := svc.tokenService.New(u.Id, 3*24*time.Hour, token.ScopeActivation)
	if err != nil {
		return nil, err
	}

	backgroundRoutinesWaitGroup.Add(1)
	background(logger, func() {
		defer backgroundRoutinesWaitGroup.Done()
		err = svc.mailerService.WelcomeNewRegisteredUser(context.Background(), u, token.Plaintext)
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
