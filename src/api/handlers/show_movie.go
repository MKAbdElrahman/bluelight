package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func newShowMovieHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parsedId, err := parseIdFromPath(r)
		if err != nil || parsedId < 1 {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "show the details of movie %d\n", parsedId)
	}
}

func parseIdFromPath(r *http.Request) (int, error) {
	idFromPath := r.PathValue("id")
	parsedId, err := strconv.Atoi(idFromPath)
	return parsedId, err
}
