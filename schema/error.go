package schema

import (
	"fmt"
	"net/http"
)

// ConnectorError represents a connector error that follows [NDC error handling specs]
//
// [NDC error handling specs]: https://hasura.github.io/ndc-spec/specification/error-handling.html
type ConnectorError struct {
	statusCode int
	// A human-readable summary of the error
	Message string `json:"message"`
	// Any additional structured information about the error
	Details map[string]any `json:"details"`
}

// StatusCode gets the inner status code
func (ce ConnectorError) StatusCode() int {
	return ce.statusCode
}

func (ce ConnectorError) String() string {
	return fmt.Sprintf("%d: %s, details: %+v", ce.statusCode, ce.Message, ce.Details)
}

func (ce ConnectorError) Error() string {
	return ce.Message
}

// BadRequestError returns an error when the request did not match the data connector's expectation based on this specification
func BadRequestError(message string, details map[string]any) *ConnectorError {
	return &ConnectorError{http.StatusBadRequest, message, details}
}

// ForbiddenError returns an error when the request could not be handled because a permission check failed,
// for example, a mutation might fail because a check constraint was not met
func ForbiddenError(message string, details map[string]any) *ConnectorError {
	return &ConnectorError{http.StatusForbidden, message, details}
}

// ConflictError returns an error when the request could not be handled because it would create a conflicting state for the data source,
// for example, a mutation might fail because a foreign key constraint was not met
func ConflictError(message string, details map[string]any) *ConnectorError {
	return &ConnectorError{http.StatusConflict, message, details}
}

// InternalServerError returns an error when the request could not be handled because of an error on the server
func InternalServerError(message string, details map[string]any) *ConnectorError {
	return &ConnectorError{http.StatusInternalServerError, message, details}
}

// NotSupportedError returns an error when the request could not be handled because it relies on an unsupported capability.
// Note: this ought to indicate an error on the caller side, since the caller should not generate requests which are incompatible with the indicated capabilities
func NotSupportedError(message string, details map[string]any) *ConnectorError {
	return &ConnectorError{http.StatusInternalServerError, message, details}
}

// UnauthorizeError returns an unauthorized error.
func UnauthorizeError(message string, details map[string]any) *ConnectorError {
	return &ConnectorError{http.StatusUnauthorized, message, details}
}

// UnprocessableContentError returns an error when the request could not be handled because, while the request was well-formed, it was not semantically correct.
// For example, a value for a custom scalar type was provided, but with an incorrect type.
func UnprocessableContentError(message string, details map[string]any) *ConnectorError {
	return &ConnectorError{http.StatusUnprocessableEntity, message, details}
}

// BadGatewayError returns an error when the request could not be handled because an upstream service was unavailable or returned an unexpected response,
// e.g., a connection to a database server failed
func BadGatewayError(message string, details map[string]any) *ConnectorError {
	return &ConnectorError{http.StatusBadGateway, message, details}
}
