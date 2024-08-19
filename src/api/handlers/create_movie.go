package handlers

import (
	"fmt"
	"net/http"
)

func newCreateMovieHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "create a new movie")
	}
}
