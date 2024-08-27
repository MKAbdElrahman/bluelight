package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bluelight.mkcodedev.com/src/lib/jsonio"
)

// UpdateMovieRequest represents the request structure for updating a movie.
type UpdateMovieRequest struct {
	Body        MovieDetails `json:"body"`
	IdPathParam int64
}

func NewUpdateMovieRequest(r *http.Request) (UpdateMovieRequest, *ClientError) {

	parsedId, err := parseIdFromPath(r)
	if err != nil {
		return UpdateMovieRequest{}, BadRequestError
	}

	var body MovieDetails

	err = jsonio.NewJSONReader().ReadJSON(r, &body)
	if err != nil {
		return UpdateMovieRequest{}, BadRequestError
	}

	req := UpdateMovieRequest{
		Body:        body,
		IdPathParam: parsedId,
	}

	vErrs := req.Validate()
	if len(vErrs.Errors) != 0 {
		return UpdateMovieRequest{}, BadRequestError.WithDetails(vErrs.Errors)
	}
	return req, nil
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

func parseIdFromPath(r *http.Request) (int64, error) {
	idFromPath := r.PathValue("id")
	parsedId, err := strconv.ParseInt(idFromPath, 10, 64)
	return parsedId, err
}
