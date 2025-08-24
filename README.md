# HTTP Error

A Go package that allows HTTP handlers to return errors instead of manually writing status codes and responses.

## Usage

```go
package main

import (
    "net/http"
    "github.com/perbu/httperror"
)

func getUser(w http.ResponseWriter, r *http.Request) error {
    user, err := findUser(r.URL.Path)
    if err != nil {
        return httperror.NotFound("User not found")
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(user)
    return nil
}

func main() {
    mux := http.NewServeMux()
    mux.Handle("/users/", httperror.NewHandler(getUser))
    http.ListenAndServe(":8080", mux)
}
```

## Error Types

```go
httperror.BadRequest("Invalid input")
httperror.Unauthorized("Authentication required")
httperror.Forbidden("Access denied")
httperror.NotFound("Resource not found")
httperror.MethodNotAllowed("Method not allowed")
httperror.Conflict("Resource conflict")
httperror.UnprocessableEntity("Invalid data")
httperror.InternalServerError("Server error")
httperror.NotImplemented("Not implemented")
httperror.ServiceUnavailable("Service unavailable")
```

## Response Formats

The package formats errors based on the `Accept` header:

- `application/json`: `{"error":"Not found","status":404,"code":"Not Found"}`
- `application/xml`: XML error response with message, status, and code
- `text/html`: Basic HTML error page
- `text/plain`: Plain text error message

## Context Support

```go
func handler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
    // Handler implementation
    return nil
}

mux.Handle("/path", httperror.NewContextHandler(handler))
```

## Custom Formatters

### Built-in Content Negotiation
```go
formatter := httperror.NewContentNegotiatingFormatter()
mux.Handle("/api", httperror.NewHandlerWithFormatter(handler, formatter))
```

### Pluggable Formatter Registration
```go
negotiator := httperror.NewContentNegotiator().
    Register("application/json", &httperror.JSONFormatter{PrettyPrint: true}).
    Register("text/html", httperror.NewHTMLFormatter()).
    Register("application/xml", &httperror.XMLFormatter{}).
    SetDefault(&httperror.TextFormatter{})

mux.Handle("/api", httperror.NewHandlerWithFormatter(handler, negotiator))
```

### Available Formatters
- `JSONFormatter` - JSON responses
- `HTMLFormatter` - HTML error pages  
- `TextFormatter` - Plain text responses
- `XMLFormatter` - XML responses