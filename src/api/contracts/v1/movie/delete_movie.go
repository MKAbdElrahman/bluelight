package v1movie

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/contracts/webutil"
)

// DeleteMovieRequest represents the request structure for deleting a movie.
type DeleteMovieRequest struct {
	IdPathParam int64
}

func NewDeleteMovieRequest(r *http.Request) (DeleteMovieRequest, *apierror.ClientError) {

	parsedId, err := webutil.ParseIdFromPath(r)
	if err != nil {
		return DeleteMovieRequest{}, apierror.BadRequestError
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
