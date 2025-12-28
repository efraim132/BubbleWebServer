package server

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test JSON Response ==================================================================================================

// TestJSON tests sending JSON responses
func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"message": "hello",
		"status":  "ok",
	}

	err := JSON(w, 200, data)
	if err != nil {
		t.Fatalf("JSON() returned error: %v", err)
	}

	// Check status code
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Check body
	var result map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["message"] != "hello" {
		t.Errorf("Expected message 'hello', got '%s'", result["message"])
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result["status"])
	}
}

// TestJSONWithDifferentTypes tests JSON with various data types
func TestJSONWithDifferentTypes(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{"string", "hello"},
		{"int", 42},
		{"bool", true},
		{"slice", []string{"a", "b", "c"}},
		{"map", map[string]int{"a": 1, "b": 2}},
		{"struct", struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{"Alice", 30}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := JSON(w, 200, tt.data)
			if err != nil {
				t.Errorf("JSON() returned error for %s: %v", tt.name, err)
			}

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		})
	}
}

// Test Error Response =================================================================================================

// TestError tests sending error responses
func TestError(t *testing.T) {
	w := httptest.NewRecorder()

	err := Error(w, 404, "not found")
	if err != nil {
		t.Fatalf("Error() returned error: %v", err)
	}

	// Check status code
	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Check body
	var result map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["error"] != "not found" {
		t.Errorf("Expected error 'not found', got '%s'", result["error"])
	}
}

// TestErrorWithCode tests sending error responses with error codes
func TestErrorWithCode(t *testing.T) {
	w := httptest.NewRecorder()

	err := ErrorWithCode(w, 400, "invalid input", "INVALID_INPUT")
	if err != nil {
		t.Fatalf("ErrorWithCode() returned error: %v", err)
	}

	// Check body
	var result map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["error"] != "invalid input" {
		t.Errorf("Expected error 'invalid input', got '%s'", result["error"])
	}

	if result["code"] != "INVALID_INPUT" {
		t.Errorf("Expected code 'INVALID_INPUT', got '%s'", result["code"])
	}
}

// Test Text Response ==================================================================================================

// TestText tests sending text responses
func TestText(t *testing.T) {
	w := httptest.NewRecorder()

	err := Text(w, 200, "Hello, World!")
	if err != nil {
		t.Fatalf("Text() returned error: %v", err)
	}

	// Check status code
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		t.Errorf("Expected Content-Type text/plain, got %s", contentType)
	}

	// Check body
	body := w.Body.String()
	if body != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got '%s'", body)
	}
}

// Test HTML Response ==================================================================================================

// TestHTML tests sending HTML responses
func TestHTML(t *testing.T) {
	w := httptest.NewRecorder()

	html := "<h1>Hello</h1>"
	err := HTML(w, 200, html)
	if err != nil {
		t.Fatalf("HTML() returned error: %v", err)
	}

	// Check status code
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		t.Errorf("Expected Content-Type text/html, got %s", contentType)
	}

	// Check body
	body := w.Body.String()
	if body != html {
		t.Errorf("Expected body '%s', got '%s'", html, body)
	}
}

// Test NoContent Response =============================================================================================

// TestNoContent tests sending 204 No Content responses
func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()

	err := NoContent(w)
	if err != nil {
		t.Fatalf("NoContent() returned error: %v", err)
	}

	// Check status code
	if w.Code != 204 {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Body should be empty
	if w.Body.Len() != 0 {
		t.Errorf("Expected empty body, got %d bytes", w.Body.Len())
	}
}

// Test Redirect =======================================================================================================

// TestRedirect tests redirects
func TestRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/old", nil)

	Redirect(w, r, "/new", 302)

	// Check status code
	if w.Code != 302 {
		t.Errorf("Expected status 302, got %d", w.Code)
	}

	// Check location header
	location := w.Header().Get("Location")
	if location != "/new" {
		t.Errorf("Expected Location /new, got %s", location)
	}
}

// Test Convenience Functions ==========================================================================================

// TestCreated tests Created() helper
func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{"id": "123"}
	err := Created(w, data)
	if err != nil {
		t.Fatalf("Created() returned error: %v", err)
	}

	if w.Code != 201 {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

// TestOK tests OK() helper
func TestOK(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{"status": "ok"}
	err := OK(w, data)
	if err != nil {
		t.Fatalf("OK() returned error: %v", err)
	}

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestBadRequest tests BadRequest() helper
func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()

	err := BadRequest(w, "invalid input")
	if err != nil {
		t.Fatalf("BadRequest() returned error: %v", err)
	}

	if w.Code != 400 {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestUnauthorized tests Unauthorized() helper
func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()

	err := Unauthorized(w, "please login")
	if err != nil {
		t.Fatalf("Unauthorized() returned error: %v", err)
	}

	if w.Code != 401 {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// TestForbidden tests Forbidden() helper
func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()

	err := Forbidden(w, "access denied")
	if err != nil {
		t.Fatalf("Forbidden() returned error: %v", err)
	}

	if w.Code != 403 {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

// TestNotFound tests NotFound() helper
func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	err := NotFound(w, "resource not found")
	if err != nil {
		t.Fatalf("NotFound() returned error: %v", err)
	}

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestServerError tests ServerError() helper
func TestServerError(t *testing.T) {
	w := httptest.NewRecorder()

	err := ServerError(w, "database error")
	if err != nil {
		t.Fatalf("ServerError() returned error: %v", err)
	}

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// Test Status Codes ===================================================================================================

// TestVariousStatusCodes tests various status codes
func TestVariousStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"200 OK", 200},
		{"201 Created", 201},
		{"400 Bad Request", 400},
		{"401 Unauthorized", 401},
		{"403 Forbidden", 403},
		{"404 Not Found", 404},
		{"500 Internal Server Error", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JSON(w, tt.statusCode, map[string]string{"test": "data"})

			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

// Benchmark Tests =====================================================================================================

// BenchmarkJSON benchmarks JSON responses
func BenchmarkJSON(b *testing.B) {
	data := map[string]string{"message": "hello", "status": "ok"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		JSON(w, 200, data)
	}
}

// BenchmarkText benchmarks text responses
func BenchmarkText(b *testing.B) {
	text := "Hello, World!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		Text(w, 200, text)
	}
}

// BenchmarkError benchmarks error responses
func BenchmarkError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		Error(w, 404, "not found")
	}
}
