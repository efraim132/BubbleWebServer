package server

import "fmt"

// Error Handling System ===============================================================================================
// This file provides a structured way to handle errors in your HTTP handlers.
// Instead of manually writing error responses every time something goes wrong,
// you can return an APIError and the server will automatically convert it to a proper HTTP response.

// APIError Type =======================================================================================================

// APIError represents an HTTP error with a status code, message, and optional error code
// This implements Go's error interface, so it can be used anywhere a normal error is expected.
//
// Fields:
//
//	Status  - HTTP status code (400, 404, 500, etc.) - NOT included in JSON response (json:"-")
//	Message - Human-readable error message - included in JSON as "error"
//	Code    - Machine-readable error code - optional, included in JSON as "code" if present
//
// Example JSON response: {"error":"User not found","code":"NOT_FOUND"}
type APIError struct {
	Status  int    `json:"-"`              // HTTP status code - excluded from JSON with json:"-"
	Message string `json:"error"`          // Error message - shows as "error" in JSON
	Code    string `json:"code,omitempty"` // Error code - only included if not empty (omitempty)
}

// Error implements the error interface for APIError
// This method is required to make APIError work as a standard Go error.
// When you print an APIError or log it, this method determines what gets displayed.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] %s", e.Code, e.Message) // Format: [NOT_FOUND] User not found
	}
	return e.Message // Format: User not found
}

// Constructor Functions ===============================================================================================
// These functions create new APIError instances

// NewError creates a new APIError with the given status and message
// Use this when you want a simple error without an error code.
//
// Example: return server.NewError(400, "Invalid email format")
// Result:  HTTP 400 with body: {"error":"Invalid email format"}
func NewError(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

// NewErrorWithCode creates a new APIError with status, message, and code
// Use this when you want to include a machine-readable error code for programmatic handling.
//
// Example: return server.NewErrorWithCode(400, "Invalid email format", "INVALID_EMAIL")
// Result:  HTTP 400 with body: {"error":"Invalid email format","code":"INVALID_EMAIL"}
func NewErrorWithCode(status int, message, code string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
		Code:    code,
	}
}

// Predefined Common Errors ============================================================================================
// These are ready-to-use errors for common scenarios.
// Just return one of these from your handler and the server will handle the rest.

// 4xx Client Errors ---------------------------------------------------------------------------------------------------
// These indicate something was wrong with the client's request

// ErrBadRequest represents a 400 Bad Request error
// Use when: The request is malformed, invalid input, missing required fields, etc.
//
// Example: if email == "" { return server.ErrBadRequest }
var ErrBadRequest = &APIError{
	Status:  400,
	Message: "Bad request",
	Code:    "BAD_REQUEST",
}

// ErrUnauthorized represents a 401 Unauthorized error
// Use when: User needs to log in / provide authentication
//
// Example: if !isLoggedIn { return server.ErrUnauthorized }
var ErrUnauthorized = &APIError{
	Status:  401,
	Message: "Unauthorized",
	Code:    "UNAUTHORIZED",
}

// ErrForbidden represents a 403 Forbidden error
// Use when: User is logged in but doesn't have permission for this action
//
// Example: if !user.IsAdmin { return server.ErrForbidden }
var ErrForbidden = &APIError{
	Status:  403,
	Message: "Forbidden",
	Code:    "FORBIDDEN",
}

// ErrNotFound represents a 404 Not Found error
// Use when: The requested resource doesn't exist
//
// Example: if user == nil { return server.ErrNotFound }
var ErrNotFound = &APIError{
	Status:  404,
	Message: "Resource not found",
	Code:    "NOT_FOUND",
}

// ErrMethodNotAllowed represents a 405 Method Not Allowed error
// Use when: The HTTP method is wrong (e.g., used POST when only GET is allowed)
//
// Example: if r.Method != "POST" { return server.ErrMethodNotAllowed }
var ErrMethodNotAllowed = &APIError{
	Status:  405,
	Message: "Method not allowed",
	Code:    "METHOD_NOT_ALLOWED",
}

// ErrConflict represents a 409 Conflict error
// Use when: The request conflicts with current state (e.g., trying to create something that already exists)
//
// Example: if emailExists { return server.ErrConflict }
var ErrConflict = &APIError{
	Status:  409,
	Message: "Conflict",
	Code:    "CONFLICT",
}

// ErrUnprocessableEntity represents a 422 Unprocessable Entity error
// Use when: Request is valid JSON but semantically incorrect
//
// Example: if age < 0 { return server.ErrUnprocessableEntity }
var ErrUnprocessableEntity = &APIError{
	Status:  422,
	Message: "Unprocessable entity",
	Code:    "UNPROCESSABLE_ENTITY",
}

// 5xx Server Errors ---------------------------------------------------------------------------------------------------
// These indicate something went wrong on the server side

// ErrInternalServer represents a 500 Internal Server Error
// Use when: Unexpected server error, database errors, etc.
//
// Example: if err := db.Query(); err != nil { return server.ErrInternalServer }
var ErrInternalServer = &APIError{
	Status:  500,
	Message: "Internal server error",
	Code:    "INTERNAL_SERVER_ERROR",
}

// ErrNotImplemented represents a 501 Not Implemented error
// Use when: Feature exists in API but isn't implemented yet
//
// Example: return server.ErrNotImplemented // Placeholder for future feature
var ErrNotImplemented = &APIError{
	Status:  501,
	Message: "Not implemented",
	Code:    "NOT_IMPLEMENTED",
}

// ErrServiceUnavailable represents a 503 Service Unavailable error
// Use when: Server is temporarily down for maintenance or overloaded
//
// Example: if maintenanceMode { return server.ErrServiceUnavailable }
var ErrServiceUnavailable = &APIError{
	Status:  503,
	Message: "Service unavailable",
	Code:    "SERVICE_UNAVAILABLE",
}

// Custom Error Helper Functions =======================================================================================
// These functions create custom error instances with your own messages
// Use these when the predefined errors aren't specific enough

// BadRequestError creates a custom 400 error with your message
// Use when: You want to explain what specific input was invalid
//
// Example: return server.BadRequestError("Email must be in valid format")
// Result:  HTTP 400 with body: {"error":"Email must be in valid format","code":"BAD_REQUEST"}
func BadRequestError(message string) *APIError {
	return NewErrorWithCode(400, message, "BAD_REQUEST")
}

// UnauthorizedError creates a custom 401 error with your message
// Use when: You want to explain why authentication failed
//
// Example: return server.UnauthorizedError("Invalid API token")
// Result:  HTTP 401 with body: {"error":"Invalid API token","code":"UNAUTHORIZED"}
func UnauthorizedError(message string) *APIError {
	return NewErrorWithCode(401, message, "UNAUTHORIZED")
}

// ForbiddenError creates a custom 403 error with your message
// Use when: You want to explain what permission is missing
//
// Example: return server.ForbiddenError("Only admins can delete users")
// Result:  HTTP 403 with body: {"error":"Only admins can delete users","code":"FORBIDDEN"}
func ForbiddenError(message string) *APIError {
	return NewErrorWithCode(403, message, "FORBIDDEN")
}

// NotFoundError creates a custom 404 error with your message
// Use when: You want to specify what resource wasn't found
//
// Example: return server.NotFoundError("User with ID 123 not found")
// Result:  HTTP 404 with body: {"error":"User with ID 123 not found","code":"NOT_FOUND"}
func NotFoundError(message string) *APIError {
	return NewErrorWithCode(404, message, "NOT_FOUND")
}

// InternalServerError creates a custom 500 error with your message
// Use when: You want to log what went wrong (but still show generic message to user)
//
// Example: return server.InternalServerError("Failed to connect to database")
// Result:  HTTP 500 with body: {"error":"Failed to connect to database","code":"INTERNAL_SERVER_ERROR"}
func InternalServerError(message string) *APIError {
	return NewErrorWithCode(500, message, "INTERNAL_SERVER_ERROR")
}

// Usage Examples ======================================================================================================
//
// 1. Using predefined errors:
//
//    func getUser(w http.ResponseWriter, r *http.Request) error {
//        user := findUser(id)
//        if user == nil {
//            return server.ErrNotFound  // Simple! Server handles the rest
//        }
//        return server.JSON(w, 200, user)
//    }
//
// 2. Using custom error messages:
//
//    func createUser(w http.ResponseWriter, r *http.Request) error {
//        if email == "" {
//            return server.BadRequestError("Email is required")  // Custom message
//        }
//        return server.Created(w, newUser)
//    }
//
// 3. Creating fully custom errors:
//
//    func customHandler(w http.ResponseWriter, r *http.Request) error {
//        return server.NewErrorWithCode(418, "I'm a teapot", "TEAPOT")  // Any status + message + code
//    }
//
// How It Works:
// When you return an APIError from your handler, the server's wrapHandler() function (in server.go)
// catches it and automatically converts it to a proper JSON HTTP response. You never have to manually
// write error responses - just return the appropriate error and you're done!
