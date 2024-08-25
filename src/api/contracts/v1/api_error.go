package v1

// Default ServerError values
var (
	InternalServerError = ServerError{
		Code:    500,
		Message: "The server encountered an unexpected condition that prevented it from fulfilling the request.",
	}

	ServiceUnavailableError = ServerError{
		Code:    503,
		Message: "The server is currently unable to handle the request due to temporary overloading or maintenance.",
	}
)

// Default ClientError values
var (
	BadRequestError = ClientError{
		Code:    400,
		Message: "The server could not understand the request due to invalid syntax.",
	}

	UnauthorizedError = ClientError{
		Code:    401,
		Message: "The request requires user authentication.",
	}

	ForbiddenError = ClientError{
		Code:    403,
		Message: "The server understood the request, but refuses to authorize it.",
	}

	NotFoundError = ClientError{
		Code:    404,
		Message: "The server can not find the requested resource.",
	}

	MethodNotAllowedError = ClientError{
		Code:    405,
		Message: "The request method is not supported for the requested resource.",
	}

	ConflictError = ClientError{
		Code:    409,
		Message: "The request could not be completed due to a conflict with the current state of the resource.",
	}

	UnprocessableEntityError = ClientError{
		Code:    422,
		Message: "The server understands the content type of the request entity, but was unable to process the contained instructions.",
	}
)

// ServerError represents errors that occur on the server-side.
type ServerError struct {
	Code    int    `json:"code"`    // HTTP status code (e.g., 500)
	Message string `json:"message"` // Human-readable error message
}

// ClientError represents errors that occur due to client actions.
type ClientError struct {
	Code    int               `json:"code"`    // HTTP status code (e.g., 400)
	Message string            `json:"message"` // Human-readable error message
	Details map[string]string `json:"details"` // Additional details (e.g., validation errors)
}

// WithDetails adds details to the ClientError and returns the updated ClientError.
func (e ClientError) WithDetails(details map[string]string) ClientError {
	// If the ClientError already has details, merge them with the new ones
	if e.Details == nil {
		e.Details = details
	} else {
		for key, value := range details {
			e.Details[key] = value
		}
	}
	return e
}
