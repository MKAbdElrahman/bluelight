package v1

import (
	"time"
)

type MovieDetails struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

func (v *validationError) addError(field, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}
	v.Errors[field] = message
}

func validateMovieDetails(details MovieDetails) validationError {
	v := validationError{Errors: make(map[string]string)}

	// Validate Title
	if details.Title == "" {
		v.addError("title", "Title is required")
	}
	if len(details.Title) > 500 {
		v.addError("title", "Title must not be more than 500 bytes long")
	}

	// Validate Year
	if details.Year == 0 {
		v.addError("year", "Year is required")
	}
	if details.Year <= 1888 {
		v.addError("year", "Year must be greater than 1888")
	}
	if details.Year > int32(time.Now().Year()) {
		v.addError("year", "Year must not be in the future")
	}

	// Validate Runtime
	if details.Runtime == 0 {
		v.addError("runtime", "Runtime is required")
	}
	if details.Runtime < 0 {
		v.addError("runtime", "Runtime must be a positive integer")
	}

	return v
}
