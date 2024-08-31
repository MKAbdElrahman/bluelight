package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"bluelight.mkcodedev.com/src/core/domain/verrors"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserId    int64
	Expiry    time.Time
	Scope     string
}

func (t Token) ValidatePlainTextForm() *verrors.ValidationError {
	if t.Plaintext == "" {
		return &verrors.ValidationError{
			Field:   "token",
			Message: "must be provided",
		}
	}

	if len(t.Plaintext) != 26 {
		return &verrors.ValidationError{
			Field:   "token",
			Message: "must be 26 bytes long",
		}
	}

	return nil
}

type TokenRepositoty interface {
	Create(*Token) error
	DeleteAllForUser(scope string, userID int64) error
}

type tokenService struct {
	tokenRepository TokenRepositoty
}

func NewTokenService(r TokenRepositoty) *tokenService {
	return &tokenService{
		tokenRepository: r,
	}
}

func (s tokenService) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = s.tokenRepository.Create(token)
	return token, err
}

func (s tokenService) DeleteAllForUser(scope string, userID int64) error {
	return s.tokenRepository.DeleteAllForUser(scope, userID)
}

func generateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}
