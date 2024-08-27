package v1

import (
	"net/http"
)

// DeleteMovieRequest represents the request structure for deleting a movie.
type DeleteMovieRequest struct {
	IdPathParam int64
}

func NewDeleteMovieRequest(r *http.Request) (DeleteMovieRequest, *ClientError) {

	parsedId, err := parseIdFromPath(r)
	if err != nil {
		return DeleteMovieRequest{}, BadRequestError
	}

	req := DeleteMovieRequest{
		IdPathParam: parsedId,
	}

	return req, nil
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
