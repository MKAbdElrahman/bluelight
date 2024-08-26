package v1

import (
	"fmt"
	"net/http"
)

// CreateMovieRequest represents the request structure for creating a movie.
type CreateMovieRequest struct {
	Body MovieDetails `json:"body"`
}

// CreateMovieResponse represents the response structure for creating a movie.
type CreateMovieResponse struct {
	Id               int64    `json:"id"`
	Title            string   `json:"title"`
	Year             int32    `json:"year"`
	RuntimeInMinutes int32    `json:"runtime"`
	Genres           []string `json:"genres"`
	Version          int32    `json:"version"`
}

// Status returns the HTTP status code for the CreateMovieResponse.
func (r CreateMovieResponse) Status() int {
	return http.StatusCreated
}

// Headers returns the HTTP headers for the CreateMovieResponse.
func (r CreateMovieResponse) Headers() http.Header {
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", r.Id))
	return headers
}

// Validate validates the CreateMovieRequest.
func (r CreateMovieRequest) Validate() validationError {
	return validateMovieDetails(r.Body)
}
