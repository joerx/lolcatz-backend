package errors

import (
	"fmt"
	"net/http"
)

// Error represents a HTTP error type allowing us to generate a response matching
// the specific type of error
type Error struct {
	StatusCode int
	Message    string
}

func (e Error) Error() string {
	return fmt.Sprintf("http.Error: {StatusCode=%d Message=%s}", e.StatusCode, e.Message)
}

// Internal creates an Error representing an internal server error (500)
func Internal(err error) error {
	return Error{StatusCode: http.StatusInternalServerError, Message: err.Error()}
}

// BadRequest creates an Error representing a bad user request (400)
func BadRequest(msg string) error {
	return Error{StatusCode: http.StatusBadRequest, Message: msg}
}

// MethodNotAllowed creates an Error representing a method not allowed error (405)
func MethodNotAllowed(method string) error {
	return Error{StatusCode: http.StatusMethodNotAllowed, Message: fmt.Sprintf("Method not allowed: %s", method)}
}
