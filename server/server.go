package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Types and Structs ===================================================================================================

// Server represents an HTTP server with fluent configuration
// This is the main struct that holds all server state and configuration.
// It uses the "fluent API" pattern where each method returns *Server,
// allowing you to chain method calls together like: server.New().Port(8080).Route("/hello", handler).Build()
type Server struct {
	port           int                // The port the server listens on (e.g., 8080)
	host           string             // The host address (usually "0.0.0.0" for all interfaces)
	mux            *http.ServeMux     // Go's built-in HTTP request router - matches URLs to handlers
	httpServer     *http.Server       // The underlying Go HTTP server that actually handles connections
	routes         []Route            // Slice storing all registered routes before the server is built
	onRequest      func(RequestInfo)  // Optional callback function called every time a request is received
	onStart        func()             // Optional callback function called when server starts
	onStop         func()             // Optional callback function called when server stops
	shutdownCtx    context.Context    // Context used for signaling shutdown (advanced: for graceful shutdown patterns)
	shutdownCancel context.CancelFunc // Function to cancel the shutdown context
	mu             sync.Mutex         // Mutex (lock) to prevent race conditions when checking if server is running
	running        bool               // Boolean flag tracking whether the server is currently running
}

// Route represents a single HTTP route (endpoint)
// For example: GET /hello or POST /api/users
type Route struct {
	Method  string      // HTTP method (GET, POST, PUT, DELETE, etc.) - currently not enforced, for future use
	Path    string      // URL path like "/hello" or "/api/users"
	Handler HandlerFunc // The function that handles requests to this path
}

// HandlerFunc is a custom handler function type that can return errors
// This is different from Go's standard http.HandlerFunc because it returns an error.
// If your handler returns an error, the server will automatically convert it to an HTTP error response.
// Signature: func(w http.ResponseWriter, r *http.Request) error
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// RequestInfo contains information about an incoming HTTP request
// This is passed to the OnRequest callback so you can log or monitor requests
type RequestInfo struct {
	Method string    // HTTP method (GET, POST, etc.)
	Path   string    // URL path that was requested
	IP     string    // IP address of the client making the request
	Time   time.Time // Timestamp when the request was received
}

// Constructor and Builder Methods =====================================================================================

// New creates a new Server instance with sensible default settings
// This is the starting point for creating a server. Call this first, then chain other methods.
// Example: srv := server.New().Port(8080).Route("/hello", handler).Build()
func New() *Server {
	// Create a context for shutdown signaling (advanced pattern for clean shutdown)
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		port:           8080,      // Default port - can be changed with Port() method
		host:           "0.0.0.0", // Listen on all network interfaces (0.0.0.0 means "any address")
		routes:         []Route{}, // Empty slice to store routes
		shutdownCtx:    ctx,       // Context for shutdown coordination
		shutdownCancel: cancel,    // Function to trigger shutdown
	}
}

// Fluent Configuration Methods ========================================================================================
// These methods all return *Server so you can chain them together
// Example: server.New().Port(8080).Host("localhost").Route("/hello", handler)

// Port sets the server port (fluent method - returns *Server for chaining)
// The port is the number your server listens on (e.g., 8080 means http://localhost:8080)
// Valid ports are 1-65535, but ports below 1024 require admin/root privileges
func (s *Server) Port(port int) *Server {
	s.port = port
	return s // Return self to allow chaining
}

// Host sets the server host address (fluent method - returns *Server for chaining)
// Common values:
//
//	"0.0.0.0"   - Listen on all network interfaces (default) - accessible from other computers
//	"localhost" - Only accessible from this computer
//	"127.0.0.1" - Same as localhost
func (s *Server) Host(host string) *Server {
	s.host = host
	return s // Return self to allow chaining
}

// Route adds a route/endpoint to the server (fluent method - returns *Server for chaining)
// Example: srv.Route("/hello", myHandler)
// The path is the URL path (e.g., "/hello" matches http://localhost:8080/hello)
// The handler is your function that processes requests to this path
func (s *Server) Route(path string, handler HandlerFunc) *Server {
	s.routes = append(s.routes, Route{
		Path:    path,
		Handler: handler,
	})
	return s // Return self to allow chaining
}

// OnRequest sets a callback function that gets called for EVERY request (fluent method)
// Useful for logging, monitoring, or sending events to a UI (like the TUI)
// Example: srv.OnRequest(func(info server.RequestInfo) { log.Printf("%s %s", info.Method, info.Path) })
func (s *Server) OnRequest(callback func(RequestInfo)) *Server {
	s.onRequest = callback
	return s // Return self to allow chaining
}

// OnStart sets a callback function that gets called when the server starts (fluent method)
// Example: srv.OnStart(func() { log.Println("Server started!") })
func (s *Server) OnStart(callback func()) *Server {
	s.onStart = callback
	return s // Return self to allow chaining
}

// OnStop sets a callback function that gets called when the server stops (fluent method)
// Example: srv.OnStop(func() { log.Println("Cleaning up...") })
func (s *Server) OnStop(callback func()) *Server {
	s.onStop = callback
	return s // Return self to allow chaining
}

// Build finalizes the server configuration and registers all routes
// This method takes all the routes you've added and registers them with Go's HTTP router (mux).
// After calling Build(), the server is ready to start.
// Note: You don't have to call this manually - Start() will call it automatically if needed
func (s *Server) Build() *Server {
	// Create a new ServeMux (Go's built-in HTTP request router)
	// This is what matches incoming URLs to your handler functions
	s.mux = http.NewServeMux()

	// Register all routes with the mux
	// For each route we've stored, tell the mux "when you see this path, call this handler"
	for _, route := range s.routes {
		// wrapHandler converts our HandlerFunc (which returns errors) to Go's standard http.HandlerFunc
		s.mux.HandleFunc(route.Path, s.wrapHandler(route.Handler))
	}

	// Create the underlying HTTP server
	// This is Go's built-in server that actually listens for connections
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port), // Format: "0.0.0.0:8080"
		Handler: s.mux,                                // Use our mux to route requests
	}

	return s // Return self to allow chaining
}

// Internal Helper Methods =============================================================================================

// wrapHandler wraps our custom HandlerFunc (which returns errors) into Go's standard http.HandlerFunc
// This is the "magic" that allows handlers to return errors instead of manually writing error responses
// It also calls the OnRequest callback if one was set
func (s *Server) wrapHandler(handler HandlerFunc) http.HandlerFunc {
	// Return a standard http.HandlerFunc (which doesn't return errors)
	return func(w http.ResponseWriter, r *http.Request) {
		// First, call the OnRequest callback if one was set
		// This allows you to log/monitor every request
		if s.onRequest != nil {
			s.onRequest(RequestInfo{
				Method: r.Method,     // GET, POST, etc.
				Path:   r.URL.Path,   // The URL path that was requested
				IP:     r.RemoteAddr, // Client's IP address
				Time:   time.Now(),   // Current timestamp
			})
		}

		// Call the actual handler and check if it returned an error
		if err := handler(w, r); err != nil {
			// Handler returned an error - we need to convert it to an HTTP response

			// Check if it's our custom APIError type (which has a status code and message)
			if apiErr, ok := err.(*APIError); ok {
				// It's an APIError - send the status and message from the error
				Error(w, apiErr.Status, apiErr.Message)
			} else {
				// It's some other error - log it and send a generic 500 error
				log.Printf("Handler error: %v", err)
				Error(w, 500, "Internal server error")
			}
		}
		// If no error was returned, the handler already wrote its response (using JSON, Text, etc.)
	}
}

// Server Control Methods ==============================================================================================

// Start starts the HTTP server in a goroutine (non-blocking)
// After calling Start(), the server runs in the background and you can continue doing other things.
// Returns an error if the server is already running or if there's a startup problem.
func (s *Server) Start() error {
	// Lock the mutex to prevent race conditions
	// (what if two goroutines try to start the server at the same time?)
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is already running")
	}

	// Build the server if it hasn't been built yet
	// This allows you to call Start() without calling Build() first
	if s.httpServer == nil {
		s.Build()
	}

	// Update the server address in case port or host changed
	// This ensures the server uses the current port/host settings
	s.httpServer.Addr = fmt.Sprintf("%s:%d", s.host, s.port)

	s.running = true // Mark server as running
	s.mu.Unlock()    // Release the lock

	// Call the OnStart callback if one was set
	if s.onStart != nil {
		s.onStart()
	}

	log.Printf("Server starting on %s:%d", s.host, s.port)

	// Start the server in a goroutine (background thread)
	// This allows the Start() method to return immediately instead of blocking
	go func() {
		// ListenAndServe blocks until the server is shut down
		// It returns http.ErrServerClosed when shut down gracefully (which is normal)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// If it's an error OTHER than ErrServerClosed, log it (something went wrong)
			log.Printf("Server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the HTTP server
// "Gracefully" means it waits for active requests to finish (up to 5 seconds) before stopping.
// This is better than just killing the server immediately, which could cut off active requests.
func (s *Server) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is not running")
	}
	s.mu.Unlock()

	log.Println("Server shutting down...")

	// Create a context with a 5-second timeout
	// This means "try to shut down gracefully, but if it takes more than 5 seconds, force it"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Make sure we call cancel() when done (releases resources)

	// Shutdown gracefully - waits for active requests to finish (up to timeout)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	s.mu.Lock()
	s.running = false // Mark server as not running
	s.mu.Unlock()

	// Call the OnStop callback if one was set
	if s.onStop != nil {
		s.onStop()
	}

	log.Println("Server stopped")
	return nil
}

// Restart stops the server and then starts it again
// Useful when you want to apply configuration changes
func (s *Server) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond) // Brief pause to ensure everything is cleaned up
	return s.Start()
}

// Utility Methods =====================================================================================================

// IsRunning returns whether the server is currently running
// Thread-safe (uses mutex to prevent race conditions)
func (s *Server) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock() // Unlock when function returns
	return s.running
}

// GetPort returns the server's configured port
func (s *Server) GetPort() int {
	return s.port
}

// GetHost returns the server's configured host
func (s *Server) GetHost() string {
	return s.host
}

// Wait blocks until the shutdown context is cancelled
// Advanced: This is used for keeping the main program running until Shutdown() is called
func (s *Server) Wait() {
	<-s.shutdownCtx.Done() // Block until context is cancelled
}

// Shutdown cancels the shutdown context, signaling Wait() to return
// Advanced: Pair this with Wait() for graceful shutdown patterns
func (s *Server) Shutdown() {
	s.shutdownCancel()
}
