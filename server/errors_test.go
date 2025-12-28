package server

import (
	"testing"
)

// Test APIError Type ==================================================================================================

// TestAPIErrorImplementsError tests that APIError implements the error interface
func TestAPIErrorImplementsError(t *testing.T) {
	err := &APIError{
		Status:  404,
		Message: "not found",
		Code:    "NOT_FOUND",
	}

	// Should be able to use as error
	var _ error = err

	// Error() should return the message
	errorString := err.Error()
	if errorString != "[NOT_FOUND] not found" {
		t.Errorf("Expected '[NOT_FOUND] not found', got '%s'", errorString)
	}
}

// TestAPIErrorWithoutCode tests Error() without a code
func TestAPIErrorWithoutCode(t *testing.T) {
	err := &APIError{
		Status:  500,
		Message: "internal error",
	}

	errorString := err.Error()
	if errorString != "internal error" {
		t.Errorf("Expected 'internal error', got '%s'", errorString)
	}
}

// Test Constructor Functions ==========================================================================================

// TestNewError tests NewError constructor
func TestNewError(t *testing.T) {
	err := NewError(400, "bad request")

	if err.Status != 400 {
		t.Errorf("Expected status 400, got %d", err.Status)
	}

	if err.Message != "bad request" {
		t.Errorf("Expected message 'bad request', got '%s'", err.Message)
	}

	if err.Code != "" {
		t.Errorf("Expected empty code, got '%s'", err.Code)
	}
}

// TestNewErrorWithCode tests NewErrorWithCode constructor
func TestNewErrorWithCode(t *testing.T) {
	err := NewErrorWithCode(404, "not found", "NOT_FOUND")

	if err.Status != 404 {
		t.Errorf("Expected status 404, got %d", err.Status)
	}

	if err.Message != "not found" {
		t.Errorf("Expected message 'not found', got '%s'", err.Message)
	}

	if err.Code != "NOT_FOUND" {
		t.Errorf("Expected code 'NOT_FOUND', got '%s'", err.Code)
	}
}

// Test Predefined Errors ==============================================================================================

// TestPredefinedErrors tests all predefined error constants
func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    *APIError
		status int
		code   string
	}{
		{"ErrBadRequest", ErrBadRequest, 400, "BAD_REQUEST"},
		{"ErrUnauthorized", ErrUnauthorized, 401, "UNAUTHORIZED"},
		{"ErrForbidden", ErrForbidden, 403, "FORBIDDEN"},
		{"ErrNotFound", ErrNotFound, 404, "NOT_FOUND"},
		{"ErrMethodNotAllowed", ErrMethodNotAllowed, 405, "METHOD_NOT_ALLOWED"},
		{"ErrConflict", ErrConflict, 409, "CONFLICT"},
		{"ErrUnprocessableEntity", ErrUnprocessableEntity, 422, "UNPROCESSABLE_ENTITY"},
		{"ErrInternalServer", ErrInternalServer, 500, "INTERNAL_SERVER_ERROR"},
		{"ErrNotImplemented", ErrNotImplemented, 501, "NOT_IMPLEMENTED"},
		{"ErrServiceUnavailable", ErrServiceUnavailable, 503, "SERVICE_UNAVAILABLE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Status != tt.status {
				t.Errorf("%s: Expected status %d, got %d", tt.name, tt.status, tt.err.Status)
			}

			if tt.err.Code != tt.code {
				t.Errorf("%s: Expected code '%s', got '%s'", tt.name, tt.code, tt.err.Code)
			}

			if tt.err.Message == "" {
				t.Errorf("%s: Message should not be empty", tt.name)
			}
		})
	}
}

// Test Custom Error Functions =========================================================================================

// TestBadRequestError tests BadRequestError() function
func TestBadRequestError(t *testing.T) {
	err := BadRequestError("invalid email")

	if err.Status != 400 {
		t.Errorf("Expected status 400, got %d", err.Status)
	}

	if err.Message != "invalid email" {
		t.Errorf("Expected message 'invalid email', got '%s'", err.Message)
	}

	if err.Code != "BAD_REQUEST" {
		t.Errorf("Expected code 'BAD_REQUEST', got '%s'", err.Code)
	}
}

// TestUnauthorizedError tests UnauthorizedError() function
func TestUnauthorizedError(t *testing.T) {
	err := UnauthorizedError("invalid token")

	if err.Status != 401 {
		t.Errorf("Expected status 401, got %d", err.Status)
	}

	if err.Message != "invalid token" {
		t.Errorf("Expected message 'invalid token', got '%s'", err.Message)
	}

	if err.Code != "UNAUTHORIZED" {
		t.Errorf("Expected code 'UNAUTHORIZED', got '%s'", err.Code)
	}
}

// TestForbiddenError tests ForbiddenError() function
func TestForbiddenError(t *testing.T) {
	err := ForbiddenError("admin only")

	if err.Status != 403 {
		t.Errorf("Expected status 403, got %d", err.Status)
	}

	if err.Message != "admin only" {
		t.Errorf("Expected message 'admin only', got '%s'", err.Message)
	}

	if err.Code != "FORBIDDEN" {
		t.Errorf("Expected code 'FORBIDDEN', got '%s'", err.Code)
	}
}

// TestNotFoundError tests NotFoundError() function
func TestNotFoundError(t *testing.T) {
	err := NotFoundError("user not found")

	if err.Status != 404 {
		t.Errorf("Expected status 404, got %d", err.Status)
	}

	if err.Message != "user not found" {
		t.Errorf("Expected message 'user not found', got '%s'", err.Message)
	}

	if err.Code != "NOT_FOUND" {
		t.Errorf("Expected code 'NOT_FOUND', got '%s'", err.Code)
	}
}

// TestInternalServerError tests InternalServerError() function
func TestInternalServerError(t *testing.T) {
	err := InternalServerError("database error")

	if err.Status != 500 {
		t.Errorf("Expected status 500, got %d", err.Status)
	}

	if err.Message != "database error" {
		t.Errorf("Expected message 'database error', got '%s'", err.Message)
	}

	if err.Code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("Expected code 'INTERNAL_SERVER_ERROR', got '%s'", err.Code)
	}
}

// Test Error Usage in Context =========================================================================================

// TestErrorAsReturnValue tests using APIError as return value
func TestErrorAsReturnValue(t *testing.T) {
	handler := func() error {
		return ErrNotFound
	}

	err := handler()

	// Should be able to use as error
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Should be able to type assert
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Error("Expected *APIError type")
	}

	if apiErr.Status != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.Status)
	}
}

// TestErrorComparison tests comparing errors
func TestErrorComparison(t *testing.T) {
	err1 := ErrNotFound
	err2 := ErrNotFound

	// Same predefined error should be the same instance
	if err1 != err2 {
		t.Error("Expected same error instance")
	}

	// Different custom errors should be different
	custom1 := NotFoundError("user not found")
	custom2 := NotFoundError("user not found")

	if custom1 == custom2 {
		t.Error("Expected different error instances")
	}
}

// Benchmark Tests =====================================================================================================

// BenchmarkNewError benchmarks creating new errors
func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewError(404, "not found")
	}
}

// BenchmarkNewErrorWithCode benchmarks creating errors with codes
func BenchmarkNewErrorWithCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewErrorWithCode(404, "not found", "NOT_FOUND")
	}
}

// BenchmarkCustomErrorFunctions benchmarks custom error functions
func BenchmarkCustomErrorFunctions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NotFoundError("resource not found")
	}
}

// BenchmarkErrorString benchmarks Error() method
func BenchmarkErrorString(b *testing.B) {
	err := ErrNotFound

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
