package main

// RESTful API Example =====================================================================================================
// This example demonstrates how to build a RESTful Todo API using the BubbleWebServer package.
// It shows CRUD operations (Create, Read, Update, Delete), JSON request/response handling,
// method validation, and error handling.
//
// To run: go run main.go
// Then test with curl or your browser
//
// RESTful API Concepts:
//   REST = Representational State Transfer
//   It's a style of designing APIs where URLs represent resources (things) and HTTP methods represent actions:
//     GET    - Read/retrieve data
//     POST   - Create new data
//     PUT    - Update existing data (replace)
//     PATCH  - Partial update
//     DELETE - Delete data

import (
	"BubbleWebServer/server"
	"encoding/json" // For parsing JSON from request bodies
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv" // For converting strings to numbers (e.g., "123" -> 123)
	"strings" // For string manipulation (splitting, trimming, etc.)
	"syscall"
)

// Data Models ============================================================================================================

// Todo represents a single todo item
// The `json:"..."` tags tell Go how to convert between Go structs and JSON
// Example JSON: {"id":1,"title":"Learn Go","completed":true}
type Todo struct {
	ID        int    `json:"id"`        // Becomes "id" in JSON
	Title     string `json:"title"`     // Becomes "title" in JSON
	Completed bool   `json:"completed"` // Becomes "completed" in JSON
}

// In-Memory Data Store ===================================================================================================
// In a real application, this would be a database (PostgreSQL, MySQL, MongoDB, etc.)
// For this example, we're just using global variables to keep it simple

var todos = []Todo{
	{ID: 1, Title: "Learn Go", Completed: true},
	{ID: 2, Title: "Build a web server", Completed: true},
	{ID: 3, Title: "Deploy to production", Completed: false},
}

// nextID keeps track of the next ID to assign to a new todo
var nextID = 4

// Main Function ==========================================================================================================

func main() {
	// Create a RESTful API server
	// Notice we're setting up routes that follow REST conventions:
	//   /api/todos       - List all todos (GET)
	//   /api/todos       - Create a todo (POST) - using /create for simplicity
	//   /api/todos/:id   - Get/Delete a specific todo (GET/DELETE)
	srv := server.New().
		Port(3000).                                    // Run on port 3000
		Route("/api/todos", listTodosHandler).         // GET /api/todos - list all todos
		Route("/api/todos/", todoHandler).             // GET/DELETE /api/todos/:id - single todo operations
		Route("/api/todos/create", createTodoHandler). // POST /api/todos/create - create new todo
		Route("/health", healthHandler).               // GET /health - health check (common in production APIs)
		OnRequest(func(info server.RequestInfo) {      // Log every request
			log.Printf("[%s] %s", info.Method, info.Path)
		}).
		Build()

	// Print helpful information when the server starts
	log.Println("Starting Todo API server on http://localhost:3000")
	log.Println("Endpoints:")
	log.Println("  GET    /api/todos       - List all todos")
	log.Println("  POST   /api/todos/create - Create a new todo")
	log.Println("  GET    /api/todos/:id   - Get a specific todo")
	log.Println("  DELETE /api/todos/:id   - Delete a todo")
	log.Println("  GET    /health          - Health check")

	// Start the server
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Graceful shutdown - wait for Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	srv.Stop()
}

// API Handlers ===========================================================================================================

// healthHandler provides a simple health check endpoint
// This is standard in production APIs - monitoring tools can ping this to check if the server is running
// Example: GET http://localhost:3000/health
// Response: {"status":"healthy"}
func healthHandler(w http.ResponseWriter, r *http.Request) error {
	return server.JSON(w, 200, map[string]string{
		"status": "healthy",
	})
}

// listTodosHandler returns all todos
// Example: GET http://localhost:3000/api/todos
// Response: {"todos":[...],"count":3}
func listTodosHandler(w http.ResponseWriter, r *http.Request) error {
	// Validate HTTP method - this endpoint only accepts GET requests
	// r.Method is the HTTP method (GET, POST, DELETE, etc.)
	if r.Method != http.MethodGet {
		// If someone tries POST, PUT, etc. return a 405 Method Not Allowed error
		return server.ErrMethodNotAllowed
	}

	// Return all todos as JSON
	return server.JSON(w, 200, map[string]interface{}{
		"todos": todos,      // The todo array
		"count": len(todos), // Total number of todos (helpful for clients)
	})
}

// createTodoHandler creates a new todo
// Example: POST http://localhost:3000/api/todos/create
// Request body: {"title":"New todo"}
// Response: {"id":4,"title":"New todo","completed":false}
func createTodoHandler(w http.ResponseWriter, r *http.Request) error {
	// Only accept POST requests (creating something = POST)
	if r.Method != http.MethodPost {
		return server.ErrMethodNotAllowed
	}

	// Parse the JSON request body --------------------------------------------------------------------------------------
	// Define a struct for the request data we expect
	// This is anonymous struct (no name, just used here)
	var req struct {
		Title string `json:"title"` // We expect JSON like: {"title":"something"}
	}

	// json.NewDecoder(r.Body) creates a JSON decoder from the request body
	// .Decode(&req) parses the JSON and fills the req struct
	// If the JSON is invalid or malformed, this returns an error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Bad JSON - return a 400 Bad Request error
		return server.BadRequestError("Invalid JSON")
	}

	// Validate the input -----------------------------------------------------------------------------------------------
	// Check if title is empty
	if req.Title == "" {
		// Missing required field - return a 400 Bad Request error
		return server.BadRequestError("Title is required")
	}

	// Create the new todo ----------------------------------------------------------------------------------------------
	todo := Todo{
		ID:        nextID,    // Assign the next available ID
		Title:     req.Title, // Use the title from the request
		Completed: false,     // New todos start as not completed
	}
	nextID++                    // Increment ID for next todo
	todos = append(todos, todo) // Add to our todos array

	// Return the created todo with HTTP 201 Created
	// 201 is the standard "resource created" status code
	return server.Created(w, todo)
}

// todoHandler handles operations on a single todo (GET and DELETE)
// This demonstrates URL parameter extraction (getting the ID from the URL)
// Examples:
//
//	GET    http://localhost:3000/api/todos/1 - Get todo with ID 1
//	DELETE http://localhost:3000/api/todos/1 - Delete todo with ID 1
func todoHandler(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path -------------------------------------------------------------------------------------
	// Example: /api/todos/123 -> we want "123"

	// r.URL.Path is the full URL path (e.g., "/api/todos/123")
	// strings.TrimPrefix removes the "/api/todos/" part, leaving just "123"
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")

	// Validate the path
	if path == "" || path == "create" {
		// Empty path (/api/todos/) or /api/todos/create should return 404
		return server.ErrNotFound
	}

	// Convert the path string to an integer
	// strconv.Atoi converts string to int ("123" -> 123)
	// If it fails (e.g., "/api/todos/abc"), err will not be nil
	id, err := strconv.Atoi(path)
	if err != nil {
		// Invalid ID format (not a number)
		return server.BadRequestError("Invalid todo ID")
	}

	// Find the todo in our array -----------------------------------------------------------------------------------
	todoIndex := -1 // Index of the todo in the array (-1 means not found)
	var todo Todo   // The todo we found

	// Loop through all todos to find the one with matching ID
	for i, t := range todos {
		if t.ID == id {
			todoIndex = i // Remember the index (for deletion)
			todo = t      // Remember the todo (for returning)
			break         // Stop looking, we found it!
		}
	}

	// Check if we found the todo
	if todoIndex == -1 {
		// Todo with this ID doesn't exist - return 404
		return server.NotFoundError(fmt.Sprintf("Todo with ID %d not found", id))
	}

	// Handle different HTTP methods --------------------------------------------------------------------------------
	// r.Method tells us what action the client wants to perform
	switch r.Method {
	case http.MethodGet:
		// GET /api/todos/:id - Return the todo
		return server.JSON(w, 200, todo)

	case http.MethodDelete:
		// DELETE /api/todos/:id - Delete the todo

		// Remove the todo from the slice
		// This uses Go's slice manipulation:
		//   todos[:todoIndex]       = everything before the todo
		//   todos[todoIndex+1:]     = everything after the todo
		//   append(..., ...)         = concatenate them (skipping the todo at todoIndex)
		todos = append(todos[:todoIndex], todos[todoIndex+1:]...)

		// Return success message
		return server.JSON(w, 200, map[string]string{
			"message": fmt.Sprintf("Todo %d deleted", id),
		})

	default:
		// Any other method (POST, PUT, PATCH, etc.) is not allowed
		return server.ErrMethodNotAllowed
	}
}

// Testing This API ========================================================================================================
//
// You can test this API using curl (a command-line tool) or tools like Postman, Insomnia, etc.
//
// 1. List all todos:
//    curl http://localhost:3000/api/todos
//
// 2. Get a specific todo:
//    curl http://localhost:3000/api/todos/1
//
// 3. Create a new todo:
//    curl -X POST http://localhost:3000/api/todos/create \
//      -H "Content-Type: application/json" \
//      -d '{"title":"My new todo"}'
//
// 4. Delete a todo:
//    curl -X DELETE http://localhost:3000/api/todos/1
//
// 5. Health check:
//    curl http://localhost:3000/health
//
// Understanding the Responses:
//   - You'll see JSON responses with appropriate status codes
//   - If something goes wrong, you'll get error JSON like: {"error":"Todo not found","code":"NOT_FOUND"}
//   - Success responses include the requested data

// Key API Concepts Explained ==============================================================================================
//
// 1. RESTful Resource Naming:
//    - Use nouns for resources: /todos (not /getTodos)
//    - Use plural for collections: /todos (not /todo)
//    - Use IDs for specific items: /todos/123
//
// 2. HTTP Methods (Verbs):
//    - GET    - Retrieve data (safe, no side effects, can be cached)
//    - POST   - Create new resources
//    - PUT    - Replace entire resource
//    - PATCH  - Update part of resource
//    - DELETE - Remove resource
//
// 3. Status Codes:
//    - 200 OK         - Successful GET/PUT/PATCH/DELETE
//    - 201 Created    - Successful POST (resource created)
//    - 400 Bad Request - Invalid input from client
//    - 404 Not Found  - Resource doesn't exist
//    - 405 Method Not Allowed - Wrong HTTP method
//
// 4. JSON Request/Response:
//    - Request: Client sends JSON in request body (for POST/PUT/PATCH)
//    - Response: Server sends JSON back with data or error
//
// 5. CRUD Operations:
//    - Create: POST /api/todos
//    - Read:   GET /api/todos (list) or GET /api/todos/:id (single)
//    - Update: PUT /api/todos/:id (not implemented in this example)
//    - Delete: DELETE /api/todos/:id
