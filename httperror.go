package httperror

import (
	"context"
	"encoding/json"
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

// FormatterFunc allows using a function as a Formatter
type FormatterFunc func(w http.ResponseWriter, r *http.Request, err HTTPError)

// Format implements the Formatter interface
func (ff FormatterFunc) Format(w http.ResponseWriter, r *http.Request, err HTTPError) {
	ff(w, r, err)
}

// NewJSONFormatter creates a JSON formatter that can be used with handlers
// This bridges to the formatters subpackage
func NewJSONFormatter(prettyPrint bool) Formatter {
	return FormatterFunc(func(w http.ResponseWriter, r *http.Request, err HTTPError) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.StatusCode())

		response := struct {
			Error  string `json:"error"`
			Status int    `json:"status"`
			Code   string `json:"code,omitempty"`
		}{
			Error:  err.Message(),
			Status: err.StatusCode(),
			Code:   http.StatusText(err.StatusCode()),
		}

		var data []byte
		if prettyPrint {
			data, _ = json.MarshalIndent(response, "", "  ")
		} else {
			data, _ = json.Marshal(response)
		}

		w.Write(data)
	})
}

// NewHTMLFormatter creates an HTML formatter that can be used with handlers
// This bridges to the formatters subpackage
func NewHTMLFormatter() Formatter {
	return FormatterFunc(func(w http.ResponseWriter, r *http.Request, err HTTPError) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(err.StatusCode())

		// Simple HTML template
		html := `<!DOCTYPE html>
<html>
<head>
    <title>Error %d</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .error-container { max-width: 600px; margin: 0 auto; }
        .error-code { font-size: 48px; color: #e74c3c; margin-bottom: 20px; }
        .error-message { font-size: 18px; color: #333; margin-bottom: 20px; }
        .error-details { font-size: 14px; color: #666; }
    </style>
</head>
<body>
    <div class="error-container">
        <div class="error-code">%d</div>
        <div class="error-message">%s</div>
        <div class="error-details">%s</div>
    </div>
</body>
</html>`
		fmt.Fprintf(w, html, err.StatusCode(), err.StatusCode(), err.Message(), http.StatusText(err.StatusCode()))
	})
}
