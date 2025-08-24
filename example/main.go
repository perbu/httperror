package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/perbu/httperror"
)

func main() {
	// Create HTTP mux
	mux := http.NewServeMux()

	// Example 1: Default formatter (simple plain text)
	mux.Handle("/users", httperror.NewHandler(listUsers))
	mux.Handle("/users/", httperror.NewHandler(getUser))

	// Example 2: Custom JSON formatter
	jsonFormatter := &JSONFormatter{}
	mux.Handle("/users/create", httperror.NewHandlerWithFormatter(createUser, jsonFormatter))

	// Example 3: Context-based handlers with default formatter
	mux.Handle("/timeout", httperror.NewContextHandler(timeoutExample))

	fmt.Println("Starting server on :8080")
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET  /users          - List all users (plain text errors)")
	fmt.Println("  GET  /users/1        - Get specific user (plain text errors)")
	fmt.Println("  GET  /users/999      - Not found error (plain text errors)")
	fmt.Println("  GET  /users/invalid  - Bad request error (plain text errors)")
	fmt.Println("  POST /users/create   - Create user (JSON formatted errors)")
	fmt.Println("  GET  /timeout        - Timeout example (plain text errors)")
	fmt.Println("  GET  /panic          - Panic recovery demo (plain text errors)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  curl http://localhost:8080/users/999          # Plain text error")
	fmt.Println("  curl http://localhost:8080/users/create       # JSON error")
	fmt.Println("  curl http://localhost:8080/timeout            # Plain text error")

	// Add panic endpoint with default formatter
	mux.Handle("/panic", httperror.NewHandler(panicExample))

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

// JSONFormatter is a custom formatter that outputs JSON error responses
type JSONFormatter struct{}

func (f *JSONFormatter) Format(w http.ResponseWriter, r *http.Request, err httperror.HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode())

	response := struct {
		Error  string `json:"error"`
		Status int    `json:"status"`
		Code   string `json:"code"`
	}{
		Error:  err.Message(),
		Status: err.StatusCode(),
		Code:   http.StatusText(err.StatusCode()),
	}

	json.NewEncoder(w).Encode(response)
}
