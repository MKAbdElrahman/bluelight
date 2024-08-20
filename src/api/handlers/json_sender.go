package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type envelope map[string]any

func sendJSON(w http.ResponseWriter, data envelope, httpStatus int, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed to marshal data to json")
	}

	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(httpStatus)
	w.Write([]byte(js))
	return nil
}


