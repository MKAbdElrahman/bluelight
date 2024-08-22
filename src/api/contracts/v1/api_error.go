package v1

type ApiError struct {
	Code    int               `json:"code"`    // HTTP status code or custom error code
	Message string            `json:"message"` // Human-readable error message
	Details map[string]string `json:"details"` // Additional details about the error, e.g., validation errors
}
