package middleware

import (
	"net/http"
)

type middlewareFunc func(next http.Handler) http.Handler
