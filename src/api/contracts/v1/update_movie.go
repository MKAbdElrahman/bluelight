package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bluelight.mkcodedev.com/src/core/domain"
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

func NewUpdateMovieRequest(r *http.Request) (UpdateMovieRequest, *ClientError) {

	parsedId, err := parseIdFromPath(r)
	if err != nil {
		return UpdateMovieRequest{}, BadRequestError
	}

	var body UpdateMovieRequestBody

	err = jsonio.NewJSONReader().ReadJSON(r, &body)
	if err != nil {
		return UpdateMovieRequest{}, BadRequestError
	}

	req := UpdateMovieRequest{
		Body:        body,
		IdPathParam: parsedId,
	}

	m := &domain.Movie{}
	validator := domain.NewMovieValidator(m)

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
		return UpdateMovieRequest{}, UnprocessableEntityError.WithDetails(map[string]string{
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

func parseIdFromPath(r *http.Request) (int64, error) {
	idFromPath := r.PathValue("id")
	parsedId, err := strconv.ParseInt(idFromPath, 10, 64)
	return parsedId, err
}
