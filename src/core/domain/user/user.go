package user

import (
	"errors"
	"time"

	"regexp"

	"bluelight.mkcodedev.com/src/core/domain/verrors"
)

var (
	emailRX       = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	AnonymousUser = &User{}
)

type User struct {
	Id           int64
	CreatedAt    time.Time
	Name         string
	Email        string
	PasswordHash passwordHash
	Activated    bool
	Version      int
}

func NewUser(name, email, password string) (*User, error) {
	if name == "" {
		return nil, &verrors.ValidationError{Field: "Name", Message: "cannot be empty"}
	}
	if len(name) > 500 {
		return nil, &verrors.ValidationError{Field: "Name", Message: "must not be more than 500 bytes long"}
	}

	if !emailRX.MatchString(email) {
		return nil, &verrors.ValidationError{Field: "Email", Message: "must be a valid email address"}
	}

	if err := validatePlainTextPassword(password); err != nil {
		return nil, &verrors.ValidationError{Field: "Password", Message: err.Error()}
	}

	passwordHash, err := newPasswordHash(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func validatePlainTextPassword(p string) error {
	if p == "" {
		return errors.New("password is required")
	}
	if len(p) < 8 {
		return errors.New("password must be at least 8 bytes long")
	}
	if len(p) > 72 {
		return errors.New("password must not be more than 72 bytes long")
	}
	return nil
}
