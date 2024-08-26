package v1

import (
	"fmt"
	"net/http"
)

// UpdateMovieRequest represents the request structure for updating a movie.
type UpdateMovieRequest struct {
	Body MovieDetails `json:"body"`
}

// UpdateMovieResponse represents the response structure for updating a movie.
type UpdateMovieResponse struct {
	Id               int64    `json:"id"`
	Title            string   `json:"title"`
	Year             int32    `json:"year"`
	RuntimeInMinutes int32    `json:"runtime"`
	Genres           []string `json:"genres"`
	Version          int32    `json:"version"`
}

// Status returns the HTTP status code for the UpdateMovieResponse.
func (r UpdateMovieResponse) Status() int {
	return http.StatusOK
}

// Headers returns the HTTP headers for the UpdateMovieResponse.
func (r UpdateMovieResponse) Headers() http.Header {
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", r.Id))
	return headers
}

// Validate validates the UpdateMovieRequest.
func (r UpdateMovieRequest) Validate() validationError {
	return validateMovieDetails(r.Body)
}
