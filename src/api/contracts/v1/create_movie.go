package v1

import (
	"fmt"
	"net/http"

	"bluelight.mkcodedev.com/src/lib/jsonio"
)

// CreateMovieRequest represents the request structure for creating a movie.
type CreateMovieRequest struct {
	Body MovieDetails `json:"body"`
}

func NewCreateMovieRequest(r *http.Request) (CreateMovieRequest, *ClientError) {

	var body MovieDetails

	err := jsonio.NewJSONReader().ReadJSON(r, &body)
	if err != nil {
		return CreateMovieRequest{}, BadRequestError
	}

	req := CreateMovieRequest{
		Body:        body,
	}

	vErrs := req.Validate()
	if len(vErrs.Errors) != 0 {
		return CreateMovieRequest{}, BadRequestError.WithDetails(vErrs.Errors)
	}
	return req, nil
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
