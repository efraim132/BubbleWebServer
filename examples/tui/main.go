package main

// TUI Server Example ======================================================================================================
// This example shows how to use BubbleWebServer with the Terminal User Interface (TUI)
// in a new project. This is the recommended way to use the server during development.
//
// To run: go run main.go
// Then in the TUI, type: /start to start the server

import (
	"BubbleWebServer/server"
	"log"
	"net/http"
)

func main() {
	// Option 1: Use Default TUI Configuration ===========================================================================
	// This is the simplest way - just call StartWithTUI() and you get a nice TUI with defaults

	// Create your server
	srv := server.New().
		Port(8080).
		Route("/", homeHandler).
		Route("/api/hello", apiHelloHandler).
		Build()

	// Start with TUI (uses default banner and colors)
	if err := srv.StartWithTUI(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Option 2: Custom TUI Configuration ================================================================================
	// If you want to customize the banner, colors, title, etc., use StartWithTUIConfig()
	//
	// Uncomment this block to try it:
	/*
			config := server.DefaultTUIConfig()
			config.Banner = `
		 __  __         ____
		|  \/  |_   _  / ___|  ___ _ ____   _____ _ __
		| |\/| | | | | \___ \ / _ \ '__\ \ / / _ \ '__|
		| |  | | |_| |  ___) |  __/ |   \ V /  __/ |
		|_|  |_|\__, | |____/ \___|_|    \_/ \___|_|
		        |___/
		`
			config.Title = "My Custom Server"
			config.ShowDebug = true // Show debug information

			srv := server.New().
				Port(8080).
				Route("/", homeHandler).
				Route("/api/hello", apiHelloHandler).
				Build()

			if err := srv.StartWithTUIConfig(config); err != nil {
				log.Fatalf("Error: %v", err)
			}
	*/
}

// Handlers ================================================================================================================

// homeHandler handles the root path
func homeHandler(w http.ResponseWriter, r *http.Request) error {
	html := `
<!DOCTYPE html>
<html>
<head><title>TUI Server Example</title></head>
<body>
    <h1>TUI Server Example</h1>
    <p>This server is running with a Terminal User Interface!</p>
    <p>Check your terminal to see this request logged in real-time.</p>
    <p>Try: <a href="/api/hello">/api/hello</a></p>
</body>
</html>
`
	return server.HTML(w, 200, html)
}

// apiHelloHandler handles API requests
func apiHelloHandler(w http.ResponseWriter, r *http.Request) error {
	return server.JSON(w, 200, map[string]string{
		"message": "Hello from the TUI server!",
		"status":  "success",
	})
}

// How This Works ==========================================================================================================
//
// 1. You create a server with your routes (same as always)
// 2. Instead of srv.Start(), you call srv.StartWithTUI()
// 3. A Terminal UI launches with:
//    - Server status information
//    - Real-time request logging
//    - Commands to control the server
//
// TUI Commands:
//   /start   - Start the HTTP server
//   /stop    - Stop the server
//   /restart - Restart the server
//   /status  - Show server status
//   /help    - Show help
//   quit     - Exit the program
//
// Benefits of Using TUI:
//   - Visual feedback - see server status and requests in real-time
//   - Easy control - start/stop with simple commands
//   - Better debugging - all activity visible in one place
//   - Professional look - looks better than plain log statements
//
// When to Use TUI vs Regular Start:
//   - Development: Use StartWithTUI() - you get visual feedback
//   - Production: Use Start() - headless, runs in background
//   - Docker/Cloud: Use Start() - TUI needs interactive terminal
//   - Debugging: Use StartWithTUI() - see what's happening
//
// The TUI is completely optional - same server code works with or without it!
