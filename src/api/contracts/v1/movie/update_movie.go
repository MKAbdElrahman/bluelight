package v1movie

import (
	"fmt"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/contracts/webutil"
	"bluelight.mkcodedev.com/src/core/domain/movie"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

// UpdateMovieRequest represents the request structure for updating a movie.
type UpdateMovieRequest struct {
	Body        UpdateMovieRequestBody `json:"body"`
	IdPathParam int64
}

type UpdateMovieRequestBody struct {
	Title   *string  `json:"title"`
	Year    *int32   `json:"year"`
	Runtime *int32   `json:"runtime"`
	Genres  []string `json:"genres"`
}

func NewUpdateMovieRequest(r *http.Request) (UpdateMovieRequest, *apierror.ClientError) {

	parsedId, err := webutil.ParseIdFromPath(r)
	if err != nil {
		return UpdateMovieRequest{}, apierror.BadRequestError
	}

	var body UpdateMovieRequestBody

	err = jsonio.NewJSONReader().ReadJSON(r, &body)
	if err != nil {
		return UpdateMovieRequest{}, apierror.BadRequestError
	}

	req := UpdateMovieRequest{
		Body:        body,
		IdPathParam: parsedId,
	}

	m := &movie.Movie{}
	validator := movie.NewMovieValidator(m)

	if req.Body.Title != nil {
		m.Title = *req.Body.Title
		validator.ValidateTitle()
	}
	if req.Body.Year != nil {
		m.Year = *req.Body.Year
		validator.ValidateYear()
	}
	if req.Body.Runtime != nil {
		m.RuntimeInMinutes = *req.Body.Runtime
		validator.ValidateRuntimeInMinutes()
	}
	if req.Body.Genres != nil {
		m.Genres = req.Body.Genres
	}

	if err := validator.Errors(); err != nil {
		return UpdateMovieRequest{}, apierror.UnprocessableEntityError.WithDetails(map[string]string{
			"validation_error": err.Error(),
		})
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
