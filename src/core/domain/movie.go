package domain

import (
	"errors"
	"time"
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

// MovieValidator provides a fluent interface for validating Movie fields.
type MovieValidator struct {
	movie *Movie
	errs  []error
}

// NewMovieValidator creates a new MovieValidator for the given Movie.
func NewMovieValidator(movie *Movie) *MovieValidator {
	return &MovieValidator{movie: movie, errs: make([]error, 0)}
}

// ValidateTitle checks if the Title field is valid.
func (v *MovieValidator) ValidateTitle() *MovieValidator {
	if v.movie.Title == "" {
		v.errs = append(v.errs, errors.New("title is required"))
	} else if len(v.movie.Title) > 500 {
		v.errs = append(v.errs, errors.New("title must not be more than 500 bytes long"))
	}
	return v
}

// ValidateYear checks if the Year field is valid.
func (v *MovieValidator) ValidateYear() *MovieValidator {
	if v.movie.Year == 0 {
		v.errs = append(v.errs, errors.New("year is required"))
	} else if v.movie.Year <= 1888 {
		v.errs = append(v.errs, errors.New("year must be greater than 1888"))
	} else if v.movie.Year > int32(time.Now().Year()) {
		v.errs = append(v.errs, errors.New("year must not be in the future"))
	}
	return v
}

// ValidateRuntimeInMinutes checks if the RuntimeInMinutes field is valid.
func (v *MovieValidator) ValidateRuntimeInMinutes() *MovieValidator {
	if v.movie.RuntimeInMinutes == 0 {
		v.errs = append(v.errs, errors.New("runtime is required"))
	} else if v.movie.RuntimeInMinutes < 0 {
		v.errs = append(v.errs, errors.New("runtime must be a positive integer"))
	}
	return v
}

// ValidateAll validates all fields of the Movie.
func (v *MovieValidator) ValidateAll() *MovieValidator {
	return v.ValidateTitle().ValidateYear().ValidateRuntimeInMinutes()
}

// Errors returns all the validation errors encountered.
func (v *MovieValidator) Errors() error {
	if len(v.errs) == 0 {
		return nil
	}

	// Combine all errors into one error message
	errorMsg := "validation errors: "
	for i, err := range v.errs {
		if i > 0 {
			errorMsg += "; "
		}
		errorMsg += err.Error()
	}
	return errors.New(errorMsg)
}
