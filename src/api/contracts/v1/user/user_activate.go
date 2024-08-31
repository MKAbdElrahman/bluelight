package user

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
)

type ActivateUserRequest struct {
	Body ActivateUserRequestBody `json:"body"`
}

type ActivateUserRequestBody struct {
	TokenPlaintext string `json:"token"`
}
type ActivateUserResponse struct {
	Activated bool `json:"activated"`
}

func NewActivateUserRequest(r *http.Request) (ActivateUserRequest, *apierror.ClientError) {

	var body ActivateUserRequestBody

	err := jsonio.NewJSONReader().ReadJSON(r, &body)
	if err != nil {
		return ActivateUserRequest{}, apierror.BadRequestError
	}

	req := ActivateUserRequest{
		Body: body,
	}
	return req, nil
}

// Status returns the HTTP status code for the CreateMovieResponse.
func (r ActivateUserResponse) Status() int {
	return http.StatusOK
}

// Headers returns the HTTP headers for the CreateMovieResponse.
func (r ActivateUserResponse) Headers() http.Header {
	headers := make(http.Header)
	return headers
}
