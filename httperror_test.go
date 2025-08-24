package httperror

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPErrorInterface(t *testing.T) {
	err := BadRequest("test message")

	if err.StatusCode() != 400 {
		t.Errorf("Expected status code 400, got %d", err.StatusCode())
	}

	if err.Message() != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message())
	}

	if err.Error() != "test message" {
		t.Errorf("Expected error string 'test message', got '%s'", err.Error())
	}
}

func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name     string
		err      HTTPError
		expected int
	}{
		{"BadRequest", BadRequest("test"), 400},
		{"Unauthorized", Unauthorized("test"), 401},
		{"Forbidden", Forbidden("test"), 403},
		{"NotFound", NotFound("test"), 404},
		{"MethodNotAllowed", MethodNotAllowed("test"), 405},
		{"Conflict", Conflict("test"), 409},
		{"UnprocessableEntity", UnprocessableEntity("test"), 422},
		{"InternalServerError", InternalServerError("test"), 500},
		{"NotImplemented", NotImplemented("test"), 501},
		{"BadGateway", BadGateway("test"), 502},
		{"ServiceUnavailable", ServiceUnavailable("test"), 503},
		{"GatewayTimeout", GatewayTimeout("test"), 504},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.StatusCode() != tt.expected {
				t.Errorf("Expected status code %d, got %d", tt.expected, tt.err.StatusCode())
			}
		})
	}
}

func TestHandler(t *testing.T) {
	// Test successful handler
	successHandler := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
		return nil
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler := NewHandler(successHandler)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "success" {
		t.Errorf("Expected 'success', got '%s'", w.Body.String())
	}
}

func TestHandlerError(t *testing.T) {
	// Test error handler
	errorHandler := func(w http.ResponseWriter, r *http.Request) error {
		return NotFound("resource not found")
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler := NewHandler(errorHandler)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "resource not found") {
		t.Errorf("Expected error message in response, got '%s'", body)
	}
}

func TestContextHandler(t *testing.T) {
	// Test context handler
	contextHandler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		if ctx.Err() != nil {
			return New(http.StatusRequestTimeout, "timeout")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
		return nil
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler := NewContextHandler(contextHandler)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestWithHeaders(t *testing.T) {
	err := BadRequest("test error")
	headers := map[string]string{
		"X-Custom-Header": "custom-value",
		"Cache-Control":   "no-cache",
	}

	errWithHeaders := WithHeaders(err, headers)

	if errWithHeaders.StatusCode() != 400 {
		t.Errorf("Expected status code 400, got %d", errWithHeaders.StatusCode())
	}

	if errWithHeaders.Message() != "test error" {
		t.Errorf("Expected message 'test error', got '%s'", errWithHeaders.Message())
	}

	resultHeaders := errWithHeaders.Headers()
	if resultHeaders["X-Custom-Header"] != "custom-value" {
		t.Errorf("Expected custom header value, got '%s'", resultHeaders["X-Custom-Header"])
	}

	if resultHeaders["Cache-Control"] != "no-cache" {
		t.Errorf("Expected cache control header value, got '%s'", resultHeaders["Cache-Control"])
	}
}

func TestAsHTTPError(t *testing.T) {
	// Test with regular error
	regularErr := errors.New("test error")
	httpErr := AsHTTPError(regularErr)

	if httpErr.StatusCode() != 500 {
		t.Errorf("Expected status code 500 for regular error, got %d", httpErr.StatusCode())
	}

	// Test with existing HTTPError
	existingErr := NotFound("not found")
	httpErr2 := AsHTTPError(existingErr)

	if httpErr2.StatusCode() != 404 {
		t.Errorf("Expected status code 404 for existing HTTPError, got %d", httpErr2.StatusCode())
	}

	if httpErr2.Message() != "not found" {
		t.Errorf("Expected message 'not found', got '%s'", httpErr2.Message())
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(400, "Bad request", originalErr)

	if wrappedErr.StatusCode() != 400 {
		t.Errorf("Expected status code 400, got %d", wrappedErr.StatusCode())
	}

	if wrappedErr.Message() != "Bad request" {
		t.Errorf("Expected message 'Bad request', got '%s'", wrappedErr.Message())
	}

	// Test unwrapping
	if basic, ok := wrappedErr.(*basicError); ok {
		if basic.Unwrap() == nil {
			t.Error("Expected wrapped error to be unwrappable")
		}
	} else {
		t.Error("Expected basicError type")
	}
}
