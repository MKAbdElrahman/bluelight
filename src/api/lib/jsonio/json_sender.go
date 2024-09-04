package jsonio

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Envelope map[string]any

func SendJSON(w http.ResponseWriter, data any, httpStatus int, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed to marshal data to json")
	}

	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(httpStatus)
	_, err = w.Write([]byte(js))
	return err
}
