package v1

type ValidationError struct {
	Errors map[string]string `json:"validation_errors"`
}

func (v ValidationError) Length() int {
	return len(v.Errors)
}
