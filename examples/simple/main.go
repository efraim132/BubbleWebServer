package main

// Simple Web Server Example ==============================================================================================
// This example demonstrates how to create a simple web server using the BubbleWebServer package.
// It shows different types of responses (HTML, Text, JSON) and error handling.
//
// To run: go run main.go
// Then visit: http://localhost:8080 in your browser

import (
	"BubbleWebServer/server" // Import our server package
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Main Function ==========================================================================================================

func main() {
	// Create and configure the server using the fluent API
	// Each method returns the server, allowing us to chain them together
	srv := server.New(). // Create new server with defaults (port 8080, host 0.0.0.0)
				Port(8080).                               // Set port to 8080 (http://localhost:8080)
				Route("/", homeHandler).                  // Register route: GET / -> homeHandler (home page)
				Route("/hello", helloHandler).            // Register route: GET /hello -> helloHandler (text response)
				Route("/api/users", usersHandler).        // Register route: GET /api/users -> usersHandler (JSON response)
				Route("/api/error", errorHandler).        // Register route: GET /api/error -> errorHandler (error demo)
				OnRequest(func(info server.RequestInfo) { // Set callback that runs for EVERY request
			// Log each request with method, path, and IP address
			log.Printf("[%s] %s from %s", info.Method, info.Path, info.IP)
		}).
		OnStart(func() { // Set callback that runs when server starts
			log.Println("Server started successfully!")
			log.Println("Try these endpoints:")
			log.Println("  http://localhost:8080/")
			log.Println("  http://localhost:8080/hello")
			log.Println("  http://localhost:8080/api/users")
			log.Println("  http://localhost:8080/api/error")
		}).
		OnStop(func() { // Set callback that runs when server stops
			log.Println("Server stopped gracefully")
		}).
		Build() // Build the server (registers all routes)

	// Start the server in a goroutine (background)
	// This allows the main function to continue and listen for shutdown signals
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Graceful Shutdown Setup --------------------------------------------------------------------------------------------
	// This section sets up a way to shut down the server cleanly when you press Ctrl+C

	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)

	// Tell the OS to send SIGINT (Ctrl+C) and SIGTERM (kill) signals to our channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block here until we receive a signal (waits for Ctrl+C)
	<-quit

	// When we get here, the user pressed Ctrl+C
	log.Println("Shutting down server...")

	// Stop the server gracefully (waits for active requests to finish)
	if err := srv.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

// HTTP Handlers ===========================================================================================================
// These functions handle requests to different URLs (routes/endpoints)
// Each handler receives:
//   - w (ResponseWriter): Used to write the HTTP response back to the client
//   - r (*Request): Contains information about the incoming request (URL, headers, body, etc.)
// Each handler returns an error (or nil if no error)

// homeHandler handles requests to the root path "/"
// This demonstrates sending an HTML response (a web page)
func homeHandler(w http.ResponseWriter, r *http.Request) error {
	// Define the HTML content as a string
	// This is a complete HTML page with styles
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>BubbleWebServer Example</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>Welcome to BubbleWebServer!</h1>
    <p>This is a simple example demonstrating the server package.</p>

    <h2>Available Endpoints:</h2>
    <div class="endpoint"><strong>GET /</strong> - This page</div>
    <div class="endpoint"><strong>GET /hello</strong> - Simple text response</div>
    <div class="endpoint"><strong>GET /api/users</strong> - JSON response</div>
    <div class="endpoint"><strong>GET /api/error</strong> - Error handling example</div>
</body>
</html>
`
	// Send the HTML response
	// server.HTML() sets Content-Type to "text/html" and sends the HTML
	// Parameters: ResponseWriter, Status Code (200 = OK), HTML content
	return server.HTML(w, 200, html)
}

// helloHandler handles requests to "/hello"
// This demonstrates sending a simple text response (not HTML, not JSON)
func helloHandler(w http.ResponseWriter, r *http.Request) error {
	// Send a plain text response
	// server.Text() sets Content-Type to "text/plain" and sends the text
	// Parameters: ResponseWriter, Status Code (200 = OK), Text content
	return server.Text(w, 200, "Hello from BubbleWebServer!\n")
}

// usersHandler handles requests to "/api/users"
// This demonstrates sending a JSON response (typical for APIs)
func usersHandler(w http.ResponseWriter, r *http.Request) error {
	// Create some example user data
	// In a real app, this would come from a database
	// map[string]interface{} means "a map with string keys and any type of value"
	users := []map[string]interface{}{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
		{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
	}

	// Send the JSON response
	// server.JSON() automatically:
	//   1. Sets Content-Type to "application/json"
	//   2. Converts the Go data structure to JSON format
	//   3. Sends it as the response body
	// Parameters: ResponseWriter, Status Code (200 = OK), Data to convert to JSON
	return server.JSON(w, 200, map[string]interface{}{
		"success": true,       // Indicates the request was successful
		"count":   len(users), // Number of users returned
		"users":   users,      // The actual user data
	})

	// The response will look like:
	// {
	//   "success": true,
	//   "count": 3,
	//   "users": [
	//     {"id": 1, "name": "Alice", "email": "alice@example.com"},
	//     {"id": 2, "name": "Bob", "email": "bob@example.com"},
	//     {"id": 3, "name": "Charlie", "email": "charlie@example.com"}
	//   ]
	// }
}

// errorHandler handles requests to "/api/error"
// This demonstrates automatic error handling - when you return an error,
// the server automatically converts it to a proper HTTP error response
func errorHandler(w http.ResponseWriter, r *http.Request) error {
	// Instead of manually creating an error response, we just return an error
	// The server's wrapHandler() function will catch this error and automatically:
	//   1. Convert it to a JSON error response
	//   2. Set the appropriate status code (404 in this case)
	//   3. Send it back to the client
	//
	// This makes error handling much simpler and more consistent!
	return server.NotFoundError("The requested resource was not found")

	// The response will be:
	// HTTP 404 Not Found
	// {
	//   "error": "The requested resource was not found",
	//   "code": "NOT_FOUND"
	// }
}

// Key Concepts Explained ==================================================================================================
//
// 1. HTTP Request/Response:
//    - A client (browser, curl, etc.) sends a REQUEST to the server (e.g., "GET /hello")
//    - The server processes the request and sends back a RESPONSE (status code + headers + body)
//
// 2. Status Codes:
//    - 200: OK - Everything worked
//    - 404: Not Found - The requested resource doesn't exist
//    - 500: Internal Server Error - Something went wrong on the server
//
// 3. Content Types:
//    - text/html: HTML web pages
//    - text/plain: Plain text
//    - application/json: JSON data (common for APIs)
//
// 4. Routes/Endpoints:
//    - A route is a URL path that maps to a handler function
//    - Example: "/hello" -> helloHandler
//    - When a request comes in for "/hello", the server calls helloHandler()
//
// 5. Handler Functions:
//    - Functions that process HTTP requests and generate responses
//    - Receive request info (w, r) and return an error (or nil)
//    - Use helper functions (JSON, Text, HTML, etc.) to send responses
//
// 6. Graceful Shutdown:
//    - Instead of killing the server immediately (which might cut off active requests),
//      graceful shutdown waits for active requests to finish before stopping
//    - This is why we use signal.Notify() and srv.Stop() instead of just exiting
