package user

import (
	"net/http"
	"time"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
)

type CreateAuthRequest struct {
	Body CreateAuthRequestBody `json:"body"`
}

type CreateAuthRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewCreateAuthTokenRequest(r *http.Request) (CreateAuthRequest, *apierror.ClientError) {
	var body CreateAuthRequestBody

	err := jsonio.NewJSONReader().ReadJSON(r, &body)
	if err != nil {
		return CreateAuthRequest{}, apierror.BadRequestError
	}

	req := CreateAuthRequest{
		Body: body,
	}
	return req, nil
}

type CreateAuthResponse struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

// Status returns the HTTP status code for the CreateMovieResponse.
func (r CreateAuthResponse) Status() int {
	return http.StatusCreated
}

// Headers returns the HTTP headers for the CreateMovieResponse.
func (r CreateAuthResponse) Headers() http.Header {
	headers := make(http.Header)
	return headers
}
