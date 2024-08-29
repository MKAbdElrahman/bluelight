package v1movie

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/contracts/webutil"
)

type ShowMovieRequest struct {
	IdPathParam int64
}

func NewShowMovieRequest(r *http.Request) (ShowMovieRequest, *apierror.ClientError) {

	parsedId, err := webutil.ParseIdFromPath(r)
	if err != nil {
		return ShowMovieRequest{}, apierror.BadRequestError
	}

	req := ShowMovieRequest{
		IdPathParam: parsedId,
	}

	return req, nil
}

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
