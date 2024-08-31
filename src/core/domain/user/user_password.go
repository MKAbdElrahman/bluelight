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
