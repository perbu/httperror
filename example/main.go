package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/perbu/httperror"
)

// User represents a user in our example
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Simple in-memory user store
var users = map[int]*User{
	1: {ID: 1, Name: "Alice", Age: 30},
	2: {ID: 2, Name: "Bob", Age: 25},
}

func main() {
	// Create HTTP mux
	mux := http.NewServeMux()

	// Example 1: Default formatter (simple content negotiation)
	mux.Handle("/users", httperror.NewHandler(listUsers))
	mux.Handle("/users/", httperror.NewHandler(getUser))

	// Example 2: Using built-in content negotiating formatter
	formatter := httperror.NewContentNegotiatingFormatter()
	mux.Handle("/users/create", httperror.NewHandlerWithFormatter(createUser, formatter))

	// Example 3: Custom pluggable formatter with XML support
	customFormatter := httperror.NewContentNegotiator().
		Register("application/json", &httperror.JSONFormatter{PrettyPrint: true}).
		Register("text/html", httperror.NewHTMLFormatter()).
		Register("application/xml", &httperror.XMLFormatter{}).
		SetDefault(&httperror.TextFormatter{})

	// Example 4: Context-based handlers with custom formatter
	mux.Handle("/timeout", httperror.NewContextHandlerWithFormatter(timeoutExample, customFormatter))

	fmt.Println("Starting server on :8080")
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET  /users          - List all users (default formatter)")
	fmt.Println("  GET  /users/1        - Get specific user (default formatter)")
	fmt.Println("  GET  /users/999      - Not found error (default formatter)")
	fmt.Println("  GET  /users/invalid  - Bad request error (default formatter)")
	fmt.Println("  POST /users/create   - Create user (content negotiation)")
	fmt.Println("  GET  /timeout        - Timeout example (with XML support)")
	fmt.Println("  GET  /panic          - Panic recovery demo (with XML support)")
	fmt.Println("")
	fmt.Println("Try different Accept headers:")
	fmt.Println("  curl -H 'Accept: application/json' http://localhost:8080/timeout")
	fmt.Println("  curl -H 'Accept: application/xml' http://localhost:8080/timeout")
	fmt.Println("  curl -H 'Accept: text/html' http://localhost:8080/timeout")

	// Add panic endpoint with custom formatter
	mux.Handle("/panic", httperror.NewHandlerWithFormatter(panicExample, customFormatter))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// listUsers returns all users
func listUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return httperror.MethodNotAllowed("")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Simple JSON encoding
	fmt.Fprint(w, `{"users":[`)
	first := true
	for _, user := range users {
		if !first {
			fmt.Fprint(w, ",")
		}
		fmt.Fprintf(w, `{"id":%d,"name":"%s","age":%d}`, user.ID, user.Name, user.Age)
		first = false
	}
	fmt.Fprint(w, `]}`)

	return nil
}

// getUser returns a specific user by ID
func getUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return httperror.MethodNotAllowed("")
	}

	// Extract ID from path
	path := r.URL.Path
	idStr := path[len("/users/"):]

	if idStr == "" {
		return httperror.BadRequest("User ID is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httperror.BadRequest("Invalid user ID format")
	}

	user, exists := users[id]
	if !exists {
		return httperror.NotFound("User not found")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"id":%d,"name":"%s","age":%d}`, user.ID, user.Name, user.Age)

	return nil
}

// createUser demonstrates simple validation
func createUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return httperror.MethodNotAllowed("")
	}

	name := r.FormValue("name")
	ageStr := r.FormValue("age")

	// Simple validation
	if name == "" {
		return httperror.BadRequest("Name is required")
	}
	if ageStr == "" {
		return httperror.BadRequest("Age is required")
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return httperror.BadRequest("Age must be a number")
	}
	if age < 0 || age > 150 {
		return httperror.BadRequest("Age must be between 0 and 150")
	}

	// Create user (simplified)
	newID := len(users) + 1
	user := &User{
		ID:   newID,
		Name: name,
		Age:  age,
	}
	users[newID] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id":%d,"name":"%s","age":%d}`, user.ID, user.Name, user.Age)

	return nil
}

// timeoutExample demonstrates context handling
func timeoutExample(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Simulate long-running operation
	select {
	case <-time.After(2 * time.Second):
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Operation completed successfully")
		return nil
	case <-ctx.Done():
		return httperror.New(http.StatusRequestTimeout, "Operation timed out")
	}
}

// panicExample demonstrates panic recovery
func panicExample(w http.ResponseWriter, r *http.Request) error {
	panic("This is a simulated panic!")
}
