package handlers

import (
	"fmt"
	"net/http"
)

func newHealthCheckHandlerFunc(env, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		js := `{"status": %q, "environment": %q, "version": %q}`
		js = fmt.Sprintf(js, "available", env, version)
		w.Write([]byte(js))
	}
}
