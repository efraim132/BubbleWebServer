package server

import (
	"encoding/json"
	"net/http"
)

// Response Helpers ====================================================================================================
// These helper functions make it easy to send common HTTP responses without repeating boilerplate code.
// Instead of manually setting headers, status codes, and encoding data every time,
// you can just call one of these functions.

// Core Response Functions =============================================================================================

// JSON sends a JSON response with the given status code and data
// This is the most common response type for APIs - it automatically:
//  1. Sets the Content-Type header to "application/json"
//  2. Sets the HTTP status code (200 for OK, 404 for Not Found, etc.)
//  3. Encodes your data as JSON and sends it
//
// Example: server.JSON(w, 200, map[string]string{"message": "Hello"})
// Result:  HTTP 200 with body: {"message":"Hello"}
func JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json") // Tell the browser this is JSON
	w.WriteHeader(status)                              // Set the HTTP status code (200, 404, etc.)
	return json.NewEncoder(w).Encode(data)             // Convert data to JSON and write it
}

// Error sends a JSON error response with the given status code and message
// This is a convenience function for sending error responses in a consistent format.
// All errors will have the same JSON structure: {"error": "message"}
//
// Example: server.Error(w, 404, "User not found")
// Result:  HTTP 404 with body: {"error":"User not found"}
func Error(w http.ResponseWriter, status int, message string) error {
	return JSON(w, status, map[string]string{
		"error": message,
	})
}

// ErrorWithCode sends a JSON error response with status, message, and error code
// Similar to Error(), but includes an error code for programmatic error handling.
// Useful when clients need to handle different errors differently.
//
// Example: server.ErrorWithCode(w, 400, "Invalid email format", "INVALID_EMAIL")
// Result:  HTTP 400 with body: {"error":"Invalid email format","code":"INVALID_EMAIL"}
func ErrorWithCode(w http.ResponseWriter, status int, message, code string) error {
	return JSON(w, status, map[string]string{
		"error": message,
		"code":  code,
	})
}

// Text sends a plain text response with the given status code and text
// Use this for simple text responses (not JSON, not HTML).
//
// Example: server.Text(w, 200, "Hello, World!")
// Result:  HTTP 200 with body: Hello, World! (Content-Type: text/plain)
func Text(w http.ResponseWriter, status int, text string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // Tell browser this is plain text
	w.WriteHeader(status)                                       // Set HTTP status code
	_, err := w.Write([]byte(text))                             // Write the text as bytes
	return err
}

// HTML sends an HTML response with the given status code and HTML content
// Use this when you want to send a web page or HTML snippet.
//
// Example: server.HTML(w, 200, "<h1>Welcome!</h1>")
// Result:  HTTP 200 with body: <h1>Welcome!</h1> (Content-Type: text/html)
func HTML(w http.ResponseWriter, status int, html string) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Tell browser this is HTML
	w.WriteHeader(status)                                      // Set HTTP status code
	_, err := w.Write([]byte(html))                            // Write the HTML as bytes
	return err
}

// NoContent sends a 204 No Content response
// HTTP 204 means "success, but there's no content to return"
// Common uses: DELETE operations, updates that don't return data
//
// Example: server.NoContent(w)
// Result:  HTTP 204 with no body
func NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent) // 204 status code
	return nil
}

// Redirect sends a redirect response to the given URL
// This tells the browser to navigate to a different URL.
// Common status codes:
//
//	301 - Permanent redirect (URL has moved forever)
//	302 - Temporary redirect (URL has moved temporarily)
//	303 - See Other (redirect after POST to prevent resubmission)
//
// Example: server.Redirect(w, r, "/login", 302)
// Result:  Browser redirects to /login
func Redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	http.Redirect(w, r, url, status) // Use Go's built-in redirect
}

// Convenience Success Response Functions ==============================================================================
// These are shortcuts for common success responses - same as calling JSON() with specific status codes

// Created sends a 201 Created response with the given data
// HTTP 201 means "a new resource was successfully created"
// Commonly used in POST endpoints that create new items
//
// Example: server.Created(w, newUser)
// Result:  HTTP 201 with JSON body containing newUser
func Created(w http.ResponseWriter, data interface{}) error {
	return JSON(w, http.StatusCreated, data) // 201 status code
}

// OK sends a 200 OK response with the given data
// HTTP 200 means "request succeeded"
// This is the most common success response
//
// Example: server.OK(w, users)
// Result:  HTTP 200 with JSON body containing users
func OK(w http.ResponseWriter, data interface{}) error {
	return JSON(w, http.StatusOK, data) // 200 status code
}

// Convenience Error Response Functions ================================================================================
// These are shortcuts for common error responses - same as calling Error() with specific status codes

// BadRequest sends a 400 Bad Request error response
// HTTP 400 means "the client sent invalid data"
// Use when: invalid input, missing required fields, malformed JSON, etc.
//
// Example: server.BadRequest(w, "Email is required")
// Result:  HTTP 400 with body: {"error":"Email is required"}
func BadRequest(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusBadRequest, message) // 400 status code
}

// Unauthorized sends a 401 Unauthorized error response
// HTTP 401 means "you need to log in first"
// Use when: user is not authenticated (not logged in)
//
// Example: server.Unauthorized(w, "Please log in")
// Result:  HTTP 401 with body: {"error":"Please log in"}
func Unauthorized(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusUnauthorized, message) // 401 status code
}

// Forbidden sends a 403 Forbidden error response
// HTTP 403 means "you don't have permission to do this"
// Use when: user is logged in but doesn't have access to this resource
//
// Example: server.Forbidden(w, "Admin access required")
// Result:  HTTP 403 with body: {"error":"Admin access required"}
func Forbidden(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusForbidden, message) // 403 status code
}

// NotFound sends a 404 Not Found error response
// HTTP 404 means "the requested resource doesn't exist"
// Use when: a requested item/page/endpoint doesn't exist
//
// Example: server.NotFound(w, "User not found")
// Result:  HTTP 404 with body: {"error":"User not found"}
func NotFound(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusNotFound, message) // 404 status code
}

// ServerError sends a 500 Internal Server Error response
// HTTP 500 means "something went wrong on the server"
// Use when: database errors, unexpected conditions, bugs
//
// Example: server.ServerError(w, "Database connection failed")
// Result:  HTTP 500 with body: {"error":"Database connection failed"}
func ServerError(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusInternalServerError, message) // 500 status code
}

// HTTP Status Code Reference ==========================================================================================
// For your reference, here are common HTTP status codes and what they mean:
//
// 2xx Success:
//   200 OK              - Request succeeded
//   201 Created         - New resource was created
//   204 No Content      - Success, but no data to return
//
// 3xx Redirection:
//   301 Moved Permanently   - Resource moved to new URL forever
//   302 Found               - Resource temporarily at different URL
//   304 Not Modified        - Resource hasn't changed (caching)
//
// 4xx Client Errors (something wrong with the request):
//   400 Bad Request         - Invalid/malformed request
//   401 Unauthorized        - Need to authenticate (log in)
//   403 Forbidden           - Authenticated but not allowed
//   404 Not Found           - Resource doesn't exist
//   405 Method Not Allowed  - Wrong HTTP method (e.g., POST instead of GET)
//   409 Conflict            - Request conflicts with current state
//   422 Unprocessable Entity - Valid format but can't process
//
// 5xx Server Errors (something wrong on the server):
//   500 Internal Server Error - Generic server error
//   501 Not Implemented       - Feature not implemented yet
//   503 Service Unavailable   - Server temporarily unavailable
