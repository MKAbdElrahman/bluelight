package user

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Mailer interface {
	WelcomeNewRegisteredUser(ctx context.Context, u *User, activationToken string) error
}

type UserRepositoty interface {
	Create(u *User) error
	GetByEmail(email string) (*User, error)
	GetByToken(tokenScope, tokenPlaintext string) (*User, error)

	Update(u *User) error
}
type TokenRepository interface {
	Create(*Token) error
	DeleteAllForUser(scope string, userID int64) error
}

type UserService struct {
	userRepository  UserRepositoty
	tokenRepository TokenRepository
	mailerService   Mailer
}

func NewUserService(ur UserRepositoty, tr TokenRepository, ms Mailer) *UserService {

	return &UserService{
		userRepository:  ur,
		mailerService:   ms,
		tokenRepository: tr,
	}
}

type UserRegisterationParams struct {
	Name     string
	Email    string
	Password string
}
type UserActivationParams struct {
	TokenPlaintext string
}

func (svc *UserService) GetUserByToken(scope string, plainToken string) (*User, error) {
	u, err := svc.userRepository.GetByToken(scope, plainToken)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (svc *UserService) ActivateUser(backgroundRoutinesWaitGroup *sync.WaitGroup, logger *slog.Logger, params UserActivationParams) (*User, error) {
	var t Token
	t.Plaintext = params.TokenPlaintext
	verr := validatePlainTextPassword(t.Plaintext)
	if verr != nil {
		return nil, verr
	}

	u, err := svc.userRepository.GetByToken(ScopeActivation, params.TokenPlaintext)
	if err != nil {
		return nil, err
	}
	u.Activated = true

	err = svc.userRepository.Update(u)
	if err != nil {
		return nil, err
	}

	err = svc.DeleteAllTokensForUser(ScopeActivation, u.Id)
	if err != nil {
		return nil, err
	}
	return u, nil
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

	token, err := svc.NewUserToken(u.Id, 3*24*time.Hour, ScopeActivation)
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

func (s *UserService) NewUserToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = s.tokenRepository.Create(token)
	return token, err
}

func (s *UserService) DeleteAllTokensForUser(scope string, userID int64) error {
	return s.tokenRepository.DeleteAllForUser(scope, userID)
}

func (s UserService) CreateAuthToken(params CreateAuthTokenParams) (*Token, error) {

	if err := params.Validate(); err != nil {
		return nil, err
	}
	u, err := s.userRepository.GetByEmail(params.Email)
	if err != nil {
		return nil, err
	}
	match, err := u.PasswordHash.isHashedFrom(params.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, ErrInvalidCredentials
	}

	t, err := s.NewUserToken(u.Id, 24*time.Hour, ScopeAuthentication)
	if err != nil {
		return nil, err
	}

	return t, nil
}
