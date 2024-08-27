package v1

// validationError represents errors intended to be read by external clients
type validationError struct {
	Errors map[string]string `json:"validation_errors"`
}

// addError adds a new error message to the validationError
func (v validationError) addError(field, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}
	v.Errors[field] = message
}

func (v validationError) Length() int {
	return len(v.Errors)
}
