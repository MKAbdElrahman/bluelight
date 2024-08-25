package v1

import "time"

type CreateMovieRequest struct {
	Body struct {
		Title   string   `json:"title"`
		Year    int      `json:"year"`
		Runtime int      `json:"runtime"`
		Genres  []string `json:"genres"`
	}
}

func (r CreateMovieRequest) Validate() validationError {
	v := validationError{Errors: make(map[string]string)}

	// Validate Title
	if r.Body.Title == "" {
		v.Errors["title"] = "Title is required"
	}
	if len(r.Body.Title) > 500 {
		v.Errors["title"] = "Title must not be more than 500 bytes long"
	}

	// Validate Year
	if r.Body.Year == 0 {
		v.Errors["year"] = "Year is required"
	}

	if r.Body.Year <= 1888 { // The first movie was made in 1888
		v.Errors["year"] = "Year must be greater than 1888"
	}

	if r.Body.Year > time.Now().Year() {
		v.Errors["year"] = "Year must not be in the future"
	}

	// Validate Runtime
	if r.Body.Runtime == 0 {
		v.Errors["runtime"] = "Runtime is required"
	}
	
	if r.Body.Runtime < 0 {
		v.Errors["runtime"] = "Runtime must be a postive integer"
	}

	return v
}
