package v1movie

import (
	"fmt"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/core/domain/movie"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

// CreateMovieRequest represents the request structure for creating a movie.
type CreateMovieRequest struct {
	Body CreateMovieRequestBody `json:"body"`
}

type CreateMovieRequestBody struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

func NewCreateMovieRequest(r *http.Request) (CreateMovieRequest, *apierror.ClientError) {

	var body CreateMovieRequestBody

	err := jsonio.NewJSONReader().ReadJSON(r, &body)
	if err != nil {
		return CreateMovieRequest{}, apierror.BadRequestError
	}

	req := CreateMovieRequest{
		Body: body,
	}

	m := &movie.Movie{
		Title:            req.Body.Title,
		Year:             req.Body.Year,
		Genres:           req.Body.Genres,
		RuntimeInMinutes: req.Body.Runtime,
	}

	err = movie.NewMovieValidator(m).ValidateAll().Errors()
	if err != nil {
		return CreateMovieRequest{}, apierror.UnprocessableEntityError.WithDetails(map[string]string{
			"validation_error": err.Error(),
		})
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
