package userhandlers

import (
	"context"
	"net/http"

	"bluelight.mkcodedev.com/src/core/domain/user"
)

type contextKey string

const userContextKey = contextKey("user")

func StoreUserInContext(u *user.User, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, u)
	return r.WithContext(ctx)
}

func GetUserFromContext(r *http.Request) *user.User {
	u, ok := r.Context().Value(userContextKey).(*user.User)
	if !ok {
		panic("missing user value in request context")
	}
	return u
}
