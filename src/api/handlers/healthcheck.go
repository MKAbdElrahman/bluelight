package handlers

import (
	"fmt"
	"net/http"
)

func newHealthCheckHandlerFunc(env, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "status: available")
		fmt.Fprintf(w, "environment: %s\n", env)
		fmt.Fprintf(w, "version: %s\n", version)

	}
}
