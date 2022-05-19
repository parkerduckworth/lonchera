// Package failure is a error management library which provides a
// common interface for handling errors of different types. This
// package aims to enable easy translation between language runtime
// errors (such as JSON failures), HTTP errors, Weaviate server/client
// errors, etc.
package failure

import (
	"github.com/parkerduckworth/lonchera/log"
)

// Error is the generic error type passed amongst different packages
type Error struct {
	// StatusCode is the HTTP code for use in the server response
	StatusCode int `json:"statusCode"`

	// Message is the client-side description of the error
	Message string `json:"message"`

	// Cause is the error itself, reserved for server-side logging
	Cause error `json:"-"`
}

// NewError builds an error struct based off its attributes and logs the cause
func NewError(statusCode int, msg string, cause error) *Error {
	log.Errorf("status_code: %d, msg: %s, cause: %s", statusCode, msg, cause)

	return &Error{
		StatusCode: statusCode,
		Message:    msg,
		Cause:      cause,
	}
}
