package user

import (
	"errors"
	"time"

	"regexp"
)

var (
	emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
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

type UserValidator struct {
	user *User
	errs []error
}

func NewUserValidator(user *User) *UserValidator {
	return &UserValidator{user: user, errs: make([]error, 0)}
}

func (v *UserValidator) ValidateEmail() *UserValidator {
	if v.user.Email == "" {
		v.errs = append(v.errs, errors.New("email is required"))
	}
	if !emailRX.MatchString(v.user.Email) {
		v.errs = append(v.errs, errors.New("email must be a valid email address"))
	}
	return v
}

func (v *UserValidator) ValidateName() *UserValidator {
	if v.user.Name == "" {
		v.errs = append(v.errs, errors.New("name is required"))
	}
	if len(v.user.Name) > 500 {
		v.errs = append(v.errs, errors.New("name must not be more than 500 bytes long"))
	}
	return v
}
