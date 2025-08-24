package httperror

import (
	"context"
	"fmt"
	"net/http"
)

// HTTPError represents an HTTP error with status code and message
type HTTPError interface {
	error
	StatusCode() int
	Message() string
	Headers() map[string]string
}

// HandlerFunc is a function that returns an HTTPError instead of writing directly to ResponseWriter
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ContextHandlerFunc is a handler that receives context as first parameter
type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Formatter handles error formatting for different content types
type Formatter interface {
	Format(w http.ResponseWriter, r *http.Request, err HTTPError)
}

// basicError is a basic implementation of HTTPError
type basicError struct {
	code    int
	message string
	headers map[string]string
	cause   error
}

func (e *basicError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *basicError) StatusCode() int {
	return e.code
}

func (e *basicError) Message() string {
	return e.message
}

func (e *basicError) Headers() map[string]string {
	if e.headers == nil {
		return make(map[string]string)
	}
	return e.headers
}

func (e *basicError) Unwrap() error {
	return e.cause
}

// New creates a new HTTPError with the given status code and message
func New(code int, message string) HTTPError {
	return &basicError{
		code:    code,
		message: message,
		headers: make(map[string]string),
	}
}

// Wrap wraps an existing error with HTTP status code
func Wrap(code int, message string, err error) HTTPError {
	return &basicError{
		code:    code,
		message: message,
		headers: make(map[string]string),
		cause:   err,
	}
}

// WithHeaders adds headers to an HTTPError
func WithHeaders(err HTTPError, headers map[string]string) HTTPError {
	if be, ok := err.(*basicError); ok {
		newHeaders := make(map[string]string)
		for k, v := range be.headers {
			newHeaders[k] = v
		}
		for k, v := range headers {
			newHeaders[k] = v
		}
		return &basicError{
			code:    be.code,
			message: be.message,
			headers: newHeaders,
			cause:   be.cause,
		}
	}

	// For other implementations, create a new error
	newHeaders := make(map[string]string)
	for k, v := range headers {
		newHeaders[k] = v
	}
	for k, v := range err.Headers() {
		newHeaders[k] = v
	}

	return &basicError{
		code:    err.StatusCode(),
		message: err.Message(),
		headers: newHeaders,
	}
}

// AsHTTPError converts a regular error to HTTPError, defaulting to 500 if not already an HTTPError
func AsHTTPError(err error) HTTPError {
	if httpErr, ok := err.(HTTPError); ok {
		return httpErr
	}
	return InternalServerError("An unexpected error occurred") // security
}
