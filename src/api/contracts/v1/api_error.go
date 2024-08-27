package v1

import (
	"encoding/json"
	"fmt"
)

// Default ServerError values
var (
	InternalServerError = ServerError{
		Code:            500,
		InternalMessage: "The server encountered an unexpected condition that prevented it from fulfilling the request.",
	}

	ServiceUnavailableError = ServerError{
		Code:            503,
		InternalMessage: "The server is currently unable to handle the request due to temporary overloading or maintenance.",
	}
)

// Default ClientError values
var (
	BadRequestError = &ClientError{
		Code:              400,
		UserFacingMessage: "The server could not understand the request due to invalid syntax.",
	}

	UnauthorizedError = &ClientError{
		Code:              401,
		UserFacingMessage: "The request requires user authentication.",
	}

	ForbiddenError = &ClientError{
		Code:              403,
		UserFacingMessage: "The server understood the request, but refuses to authorize it.",
	}

	NotFoundError = &ClientError{
		Code:              404,
		UserFacingMessage: "The server can not find the requested resource.",
	}

	MethodNotAllowedError = &ClientError{
		Code:              405,
		UserFacingMessage: "The request method is not supported for the requested resource.",
	}

	ConflictError = &ClientError{
		Code:              409,
		UserFacingMessage: "The request could not be completed due to a conflict with the current state of the resource.",
	}

	UnprocessableEntityError = &ClientError{
		Code:              422,
		UserFacingMessage: "The server understands the content type of the request entity, but was unable to process the contained instructions.",
	}
)

// ServerError represents errors that occur on the server-side.
type ServerError struct {
	Code            int    `json:"code"`    // HTTP status code (e.g., 500)
	InternalMessage string `json:"message"` // Human-readable error message
}

// Error implements the error interface for ServerError
func (e *ServerError) Error() string {
	return fmt.Sprintf("ServerError: Code=%d, Message=%s", e.Code, e.InternalMessage)
}

// ClientError represents errors that occur due to client actions.
type ClientError struct {
	Code              int               `json:"code"`    // HTTP status code (e.g., 400)
	UserFacingMessage string            `json:"message"` // Human-readable error message
	UserFacingDetails map[string]string `json:"details"` // Additional details (e.g., validation errors)
}

// Error implements the error interface for ClientError
func (e *ClientError) Error() string {
	// Convert the UserFacingDetails map to a JSON string if it's not empty
	var details string
	if len(e.UserFacingDetails) > 0 {
		detailBytes, err := json.Marshal(e.UserFacingDetails)
		if err != nil {
			details = "details could not be parsed"
		} else {
			details = string(detailBytes)
		}
	}

	if details != "" {
		return fmt.Sprintf("ClientError: Code=%d, Message=%s, Details=%s", e.Code, e.UserFacingMessage, details)
	}
	return fmt.Sprintf("ClientError: Code=%d, Message=%s", e.Code, e.UserFacingMessage)
}

// WithDetails adds details to the ClientError and returns the updated ClientError.
func (e *ClientError) WithDetails(details map[string]string) *ClientError {
	// If the ClientError already has details, merge them with the new ones
	if e.UserFacingDetails == nil {
		e.UserFacingDetails = details
	} else {
		for key, value := range details {
			e.UserFacingDetails[key] = value
		}
	}
	return e
}
