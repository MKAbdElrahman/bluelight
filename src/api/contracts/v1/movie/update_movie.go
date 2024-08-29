package v1movie

import (
	"fmt"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/contracts/webutil"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
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
