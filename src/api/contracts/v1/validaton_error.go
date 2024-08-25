package v1

// Errors intended to be read by external clients

type validationError struct {
	Errors map[string]string `json:"validation_errors"`
}

func (v validationError) Length() int {
	return len(v.Errors)
}

func (v validationError) Details() map[string]string {
	return v.Errors
}
