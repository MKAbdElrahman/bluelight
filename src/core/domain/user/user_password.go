package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type passwordHash []byte

func newPasswordHash(plaintextPassword string) (passwordHash, error) {
	err := validatePlainTextPassword(plaintextPassword)
	if err != nil {
		return nil, err
	}
	cost := 12
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), cost)
	if err != nil {
		return nil, err
	}
	return hash, nil
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

func (p passwordHash) isHashedFrom(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
