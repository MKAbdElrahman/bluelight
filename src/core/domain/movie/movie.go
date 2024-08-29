package movie

import (
	"time"

	"bluelight.mkcodedev.com/src/core/domain/verrors"
)

type Movie struct {
	Id               int64
	CreatedAt        time.Time
	Title            string
	Year             int32
	RuntimeInMinutes int32
	Genres           []string
	Version          int32
}

func NewMovie(title string, year int32, runtime int32, genres []string) (*Movie, error) {
	if title == "" {
		return nil, verrors.NewValidationError("title", "title is empty")
	}
	if len(title) > 500 {
		return nil, verrors.NewValidationError(title, "title must not be more than 500 bytes long")
	}

	if year == 0 {
		return nil, verrors.NewValidationError("year", "year is 0")
	}
	if year <= 1888 {
		return nil, verrors.NewValidationError("year", "year must be greater than 1888")
	}
	if year > int32(time.Now().Year()) {
		return nil, verrors.NewValidationError("year", "year must not be in the future")
	}

	if runtime == 0 {
		return nil, verrors.NewValidationError("runtime", "runtime is required")
	}

	if runtime < 0 {
		return nil, verrors.NewValidationError("runtime", "runtime must be a positive integer")
	}

	if genres == nil {
		return nil, verrors.NewValidationError("genres", "genres is required")
	}
	if len(genres) == 0 {
		return nil, verrors.NewValidationError("genres", "must contain at least 1 genre")
	}
	if len(genres) > 5 {
		return nil, verrors.NewValidationError("genres", "must not contain more than 5 genres")
	}

	if !unique(genres) {
		return nil, verrors.NewValidationError("genres", "must not contain duplicate values")
	}

	return &Movie{
		Title:            title,
		Year:             year,
		RuntimeInMinutes: runtime,
		Genres:           genres,
	}, nil

}

func unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
