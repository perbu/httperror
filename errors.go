package httperror

import (
	"fmt"
	"net/http"
)

// Common HTTP errors

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) HTTPError {
	return New(http.StatusBadRequest, message)
}

// BadRequestf creates a 400 Bad Request error with formatting
func BadRequestf(format string, args ...interface{}) HTTPError {
	return New(http.StatusBadRequest, sprintf(format, args...))
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string) HTTPError {
	return New(http.StatusUnauthorized, message)
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string) HTTPError {
	return New(http.StatusForbidden, message)
}

// NotFound creates a 404 Not Found error
func NotFound(message string) HTTPError {
	if message == "" {
		message = "Not Found"
	}
	return New(http.StatusNotFound, message)
}

// MethodNotAllowed creates a 405 Method Not Allowed error
func MethodNotAllowed(message string) HTTPError {
	if message == "" {
		message = "Method Not Allowed"
	}
	return New(http.StatusMethodNotAllowed, message)
}

// Conflict creates a 409 Conflict error
func Conflict(message string) HTTPError {
	return New(http.StatusConflict, message)
}

// UnprocessableEntity creates a 422 Unprocessable Entity error
func UnprocessableEntity(message string) HTTPError {
	return New(http.StatusUnprocessableEntity, message)
}

// InternalServerError creates a 500 Internal Server Error
func InternalServerError(message string) HTTPError {
	if message == "" {
		message = "Internal Server Error"
	}
	return New(http.StatusInternalServerError, message)
}

// InternalServerErrorf creates a 500 Internal Server Error with formatting
func InternalServerErrorf(format string, args ...interface{}) HTTPError {
	return New(http.StatusInternalServerError, sprintf(format, args...))
}

// NotImplemented creates a 501 Not Implemented error
func NotImplemented(message string) HTTPError {
	if message == "" {
		message = "Not Implemented"
	}
	return New(http.StatusNotImplemented, message)
}

// BadGateway creates a 502 Bad Gateway error
func BadGateway(message string) HTTPError {
	if message == "" {
		message = "Bad Gateway"
	}
	return New(http.StatusBadGateway, message)
}

// ServiceUnavailable creates a 503 Service Unavailable error
func ServiceUnavailable(message string) HTTPError {
	if message == "" {
		message = "Service Unavailable"
	}
	return New(http.StatusServiceUnavailable, message)
}

// GatewayTimeout creates a 504 Gateway Timeout error
func GatewayTimeout(message string) HTTPError {
	if message == "" {
		message = "Gateway Timeout"
	}
	return New(http.StatusGatewayTimeout, message)
}

// sprintf is a helper to format strings
func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
