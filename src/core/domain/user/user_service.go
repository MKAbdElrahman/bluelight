package user

import (
	"context"
)

type Mailer interface {
	WelcomeNewRegisteredUser(ctx context.Context, recipientEmail, recipientName string) error
}

type UserService struct {
	userRepository UserRepositoty
	mailer         Mailer
}

func NewUserService(r UserRepositoty, mailer Mailer) *UserService {
	return &UserService{
		userRepository: r,
		mailer:         mailer,
	}
}

type UserRegisterationParams struct {
	Name     string
	Email    string
	Password string
}

func (svc *UserService) RegisterUser(params UserRegisterationParams) (*User, error) {
	u, err := NewUser(params.Name, params.Email, params.Password)
	if err != nil {
		return nil, err
	}

	err = svc.userRepository.Create(u)
	if err != nil {
		return nil, err
	}

	err = svc.mailer.WelcomeNewRegisteredUser(context.Background(), u.Email, u.Name)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (svc *UserService) UpdateUser(u *User) error {
	return svc.userRepository.Update(u)
}

func (svc *UserService) GetByEmail(email string) (*User, error) {
	return svc.userRepository.GetByEmail(email)
}
