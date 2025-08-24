# HTTP Error

A Go package that allows HTTP handlers to return errors instead of manually writing status codes and responses.

## Features

- Plain text error responses by default
- Built-in JSON and HTML formatters
- Custom formatter interface
- Context support for handlers
- Standard library only

## Quick Start

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

### Default Format
Errors are returned as plain text:
```
User not found
```

### JSON Format
```go
jsonFormatter := httperror.NewJSONFormatter(true) // true = pretty print
mux.Handle("/api/users/", httperror.NewHandlerWithFormatter(getUser, jsonFormatter))
```

Output:
```json
{
  "error": "User not found",
  "status": 404,
  "code": "Not Found"
}
```

### HTML Format
```go
htmlFormatter := httperror.NewHTMLFormatter()
mux.Handle("/web/users/", httperror.NewHandlerWithFormatter(getUser, htmlFormatter))
```

Returns an HTML error page.

## Context Support

```go
func handler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
    // Handler implementation
    return nil
}

mux.Handle("/path", httperror.NewContextHandler(handler))
```

## Custom Formatters

Implement the `Formatter` interface:

```go
type Formatter interface {
    Format(w http.ResponseWriter, r *http.Request, err HTTPError)
}

type MyCustomFormatter struct{}

func (f *MyCustomFormatter) Format(w http.ResponseWriter, r *http.Request, err HTTPError) {
    w.Header().Set("Content-Type", "application/custom")
    w.WriteHeader(err.StatusCode())
    // Custom formatting logic
}

customFormatter := &MyCustomFormatter{}
mux.Handle("/custom", httperror.NewHandlerWithFormatter(handler, customFormatter))
```

## Error Wrapping

```go
func handler(w http.ResponseWriter, r *http.Request) error {
    err := someOperation()
    if err != nil {
        return httperror.Wrap(500, "Operation failed", err)
    }
    return nil
}
```

## Adding Headers

```go
err := httperror.NotFound("Resource not found")
errWithHeaders := httperror.WithHeaders(err, map[string]string{
    "Cache-Control": "no-cache",
    "X-Custom-Header": "custom-value",
})
return errWithHeaders
```

## License

BSD 2-Clause