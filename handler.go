package httperror

import (
	"net/http"
)

// Handler wraps a HandlerFunc to implement http.Handler
type Handler struct {
	handler   HandlerFunc
	formatter Formatter
}

// NewHandler creates a new Handler with default formatter
func NewHandler(h HandlerFunc) *Handler {
	return &Handler{
		handler:   h,
		formatter: &DefaultFormatter{},
	}
}

// NewHandlerWithFormatter creates a new Handler with custom formatter
func NewHandlerWithFormatter(h HandlerFunc, formatter Formatter) *Handler {
	return &Handler{
		handler:   h,
		formatter: formatter,
	}
}

// ServeHTTP implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.handler(w, r)
	if err != nil {
		h.handleError(w, r, err)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	// Convert to HTTPError
	httpErr := AsHTTPError(err)

	// Set headers
	for key, value := range httpErr.Headers() {
		w.Header().Set(key, value)
	}

	// Format and write the error response
	if h.formatter != nil {
		h.formatter.Format(w, r, httpErr)
	} else {
		// Fallback to basic text response
		w.WriteHeader(httpErr.StatusCode())
		w.Write([]byte(httpErr.Message()))
	}
}

// ContextHandler wraps a ContextHandlerFunc to implement http.Handler
type ContextHandler struct {
	handler   ContextHandlerFunc
	formatter Formatter
}

// NewContextHandler creates a new ContextHandler with default formatter
func NewContextHandler(h ContextHandlerFunc) *ContextHandler {
	return &ContextHandler{
		handler:   h,
		formatter: &DefaultFormatter{},
	}
}

// NewContextHandlerWithFormatter creates a new ContextHandler with custom formatter
func NewContextHandlerWithFormatter(h ContextHandlerFunc, formatter Formatter) *ContextHandler {
	return &ContextHandler{
		handler:   h,
		formatter: formatter,
	}
}

// ServeHTTP implements http.Handler
func (h *ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.handler(r.Context(), w, r)
	if err != nil {
		h.handleError(w, r, err)
	}
}

func (h *ContextHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	// Convert to HTTPError
	httpErr := AsHTTPError(err)

	// Set headers
	for key, value := range httpErr.Headers() {
		w.Header().Set(key, value)
	}

	// Format and write the error response
	if h.formatter != nil {
		h.formatter.Format(w, r, httpErr)
	} else {
		// Fallback to basic text response
		w.WriteHeader(httpErr.StatusCode())
		w.Write([]byte(httpErr.Message()))
	}
}

// Convenience functions for creating handlers

// Handle creates a new Handler and registers it with a ServeMux
func Handle(pattern string, mux *http.ServeMux, handler HandlerFunc) {
	mux.Handle(pattern, NewHandler(handler))
}

// HandleFunc creates a new Handler and registers it with DefaultServeMux
func HandleFunc(pattern string, handler HandlerFunc) {
	http.Handle(pattern, NewHandler(handler))
}

// HandleContext creates a new ContextHandler and registers it with a ServeMux
func HandleContext(pattern string, mux *http.ServeMux, handler ContextHandlerFunc) {
	mux.Handle(pattern, NewContextHandler(handler))
}

// HandleContextFunc creates a new ContextHandler and registers it with DefaultServeMux
func HandleContextFunc(pattern string, handler ContextHandlerFunc) {
	http.Handle(pattern, NewContextHandler(handler))
}
