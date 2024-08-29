package user

import (
	"fmt"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
)

type RegisterUserRequest struct {
	Body RegisterUserRequestBody `json:"body"`
}

type RegisterUserRequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterUserResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Version   int    `json:"version"`
	Activated bool   `json:"activated"`
}

func NewRegisterUserRequest(r *http.Request) (RegisterUserRequest, *apierror.ClientError) {

	var body RegisterUserRequestBody

	err := jsonio.NewJSONReader().ReadJSON(r, &body)
	if err != nil {
		return RegisterUserRequest{}, apierror.BadRequestError
	}

	req := RegisterUserRequest{
		Body: body,
	}
	return req, nil
}

// Status returns the HTTP status code for the CreateMovieResponse.
func (r RegisterUserResponse) Status() int {
	return http.StatusCreated
}

// Headers returns the HTTP headers for the CreateMovieResponse.
func (r RegisterUserResponse) Headers() http.Header {
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", r.Id))
	return headers
}
