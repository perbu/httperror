package httperror

import "testing"

// Simple test to ensure package compiles
func TestPackageCompiles(t *testing.T) {
	err := BadRequest("test")
	if err.StatusCode() != 400 {
		t.Errorf("Expected 400, got %d", err.StatusCode())
	}
}
