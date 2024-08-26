package v1

import (
	"net/http"
)

// ShowMovieResponse represents the response structure for showing a movie.
type ShowMovieResponse struct {
	Id               int64    `json:"id"`
	Title            string   `json:"title"`
	Year             int32    `json:"year"`
	RuntimeInMinutes int32    `json:"runtime"`
	Genres           []string `json:"genres"`
	Version          int32    `json:"version"`
}

// Status returns the HTTP status code for the ShowMovieResponse.
func (r ShowMovieResponse) Status() int {
	return http.StatusOK
}

// Headers returns the HTTP headers for the ShowMovieResponse.
func (r ShowMovieResponse) Headers() http.Header {
	return make(http.Header)
}
