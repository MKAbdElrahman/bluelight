package v1

import (
	"net/http"
)

// DeleteMovieRequest represents the request structure for deleting a movie.
type DeleteMovieRequest struct {
}

// DeleteMovieResponse represents the response structure for deleting a movie.
type DeleteMovieResponse struct {
}

// Status returns the HTTP status code for the DeleteMovieResponse.
func (r DeleteMovieResponse) Status() int {
	return http.StatusNoContent

}

// Headers returns the HTTP headers for the DeleteMovieResponse.
func (r DeleteMovieResponse) Headers() http.Header {
	return make(http.Header)
}
