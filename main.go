package main

// BubbleWebServer with TUI ================================================================================================
// This is the main TUI application for BubbleWebServer
// It demonstrates how to use the server package with the built-in Terminal User Interface
//
// The TUI provides a graphical interface in your terminal for:
//   - Starting/stopping the server
//   - Viewing HTTP requests in real-time
//   - Controlling the server with commands
//
// To run: go run main.go

import (
	"BubbleWebServer/server"
	"fmt"
	"log"
	"net/http"
)

// Custom Banner ==========================================================================================================
// You can customize the banner shown in the TUI

const customBanner = ` /$$      /$$           /$$       /$$$$$$$            /$$       /$$       /$$
| $$  /$ | $$          | $$      | $$__  $$          | $$      | $$      | $$
| $$ /$$$| $$  /$$$$$$ | $$$$$$$ | $$  \\ $$ /$$   /$$| $$$$$$$ | $$$$$$$ | $$  /$$$$$$   /$$$$$$$
| $$/$$ $$ $$ /$$__  $$| $$__  $$| $$$$$$$ | $$  | $$| $$__  $$| $$__  $$| $$ /$$__  $$ /$$_____/
| $$$$_  $$$$| $$$$$$$$| $$  \\ $$| $$__  $$| $$  | $$| $$  \\ $$| $$  \\ $$| $$| $$$$$$$$|  $$$$$$
| $$$/ \\  $$$| $$_____/| $$  | $$| $$  \\ $$| $$  | $$| $$  | $$| $$  | $$| $$| $$_____/ \\____  $$
| $$/   \\  $$|  $$$$$$$| $$$$$$$/| $$$$$$$/|  $$$$$$/| $$$$$$$/| $$$$$$$/| $$|  $$$$$$$ /$$$$$$$/
|__/     \\__/ \\_______/|_______/ |_______/  \\______/ |_______/ |_______/ |__/ \\_______/|_______/
`

// Main Function ===========================================================================================================

func main() {
	// Configure the TUI
	// You can customize the banner, title, colors, etc.
	tuiConfig := server.DefaultTUIConfig()
	tuiConfig.Banner = customBanner              // Use our custom banner
	tuiConfig.Title = "WebBubbles Control Panel" // Custom title
	tuiConfig.ShowDebug = false                  // Set to true to see debug info

	// Create and configure the server
	// Add your routes here - these are the endpoints your server will respond to
	srv := server.New().
		Port(8090).                             // Server will run on port 8090
		Route("/hello", helloHandler).          // GET http://localhost:8090/hello
		Route("/header", headersHandler).       // GET http://localhost:8090/header
		Route("/api/status", apiStatusHandler). // GET http://localhost:8090/api/status
		Build()

	// Start the server with TUI
	// This will:
	//   1. Launch the Terminal User Interface
	//   2. Show the control panel
	//   3. Let you start/stop the server with commands
	//   4. Display incoming HTTP requests in real-time
	//
	// The TUI blocks here until you type 'quit'
	if err := srv.StartWithTUIConfig(tuiConfig); err != nil {
		log.Fatalf("TUI error: %v", err)
	}

	// When the TUI exits (user typed 'quit'), the program ends
	fmt.Println("Server shut down successfully!")
}

// HTTP Handlers ===========================================================================================================
// These are the functions that handle requests to different URLs

// helloHandler handles requests to /hello
// This demonstrates a simple text response
//
// Try it: http://localhost:8090/hello
func helloHandler(w http.ResponseWriter, r *http.Request) error {
	return server.Text(w, 200, "hello\n")
}

// headersHandler handles requests to /header
// This demonstrates reading and returning HTTP headers
//
// Try it: http://localhost:8090/header
func headersHandler(w http.ResponseWriter, r *http.Request) error {
	// Build a string with all the request headers
	headers := ""
	for name, values := range r.Header {
		for _, h := range values {
			headers += fmt.Sprintf("%v: %v\n", name, h)
		}
	}
	return server.Text(w, 200, headers)
}

// apiStatusHandler handles requests to /api/status
// This demonstrates a JSON API response
//
// Try it: http://localhost:8090/api/status
func apiStatusHandler(w http.ResponseWriter, r *http.Request) error {
	return server.JSON(w, 200, map[string]interface{}{
		"status":  "running",
		"version": "0.1A",
		"message": "Server is healthy",
	})
}

// How to Use the TUI ======================================================================================================
//
// When you run this program, you'll see a terminal interface with:
//   - The banner at the top
//   - A message area showing server activity
//   - A command prompt at the bottom (>)
//
// Commands you can type:
//   /start   - Start the HTTP server
//   /stop    - Stop the HTTP server
//   /restart - Restart the HTTP server
//   /status  - Show server status
//   /help    - Show help message
//   quit     - Exit the program
//
// Workflow:
//   1. Run the program: go run main.go
//   2. Type: /start
//   3. The server starts and you'll see "Server: Started on port 8090"
//   4. Open your browser and visit http://localhost:8090/hello
//   5. You'll see the request appear in the TUI in real-time!
//   6. Type: /stop to stop the server
//   7. Type: quit to exit
//
// What You'll See:
//   - When you start the server: "Server: Started on port 8090"
//   - When requests come in: "Server: Received GET request to /hello from 127.0.0.1:xxxxx"
//   - When you stop the server: "Server: Stopped"
//
// The TUI automatically shows all HTTP requests in real-time with color-coded messages!

// Adding More Routes ======================================================================================================
//
// To add more endpoints to your server, just add more Route() calls in main():
//
// Example:
//   srv := server.New().
//       Port(8090).
//       Route("/hello", helloHandler).
//       Route("/goodbye", goodbyeHandler).      // New route!
//       Route("/api/users", listUsersHandler).  // New route!
//       Build()
//
// Then create the handler functions:
//
//   func goodbyeHandler(w http.ResponseWriter, r *http.Request) error {
//       return server.Text(w, 200, "Goodbye!\n")
//   }
//
//   func listUsersHandler(w http.ResponseWriter, r *http.Request) error {
//       users := []string{"Alice", "Bob", "Charlie"}
//       return server.JSON(w, 200, map[string]interface{}{
//           "users": users,
//           "count": len(users),
//       })
//   }
//
// That's it! The TUI will automatically show requests to your new endpoints.

// Using This in New Projects ==============================================================================================
//
// To use this TUI server in a new project:
//
// 1. Import the server package:
//    import "BubbleWebServer/server"
//
// 2. Create your server with routes:
//    srv := server.New().
//        Port(8080).
//        Route("/myendpoint", myHandler).
//        Build()
//
// 3. Start with TUI:
//    srv.StartWithTUI()
//
// That's all you need! You get a full terminal UI for monitoring and controlling your server.
//
// If you want to run WITHOUT the TUI (for example, in production or as a background service):
//    srv.Start()  // Regular start without TUI
//
// The choice is yours - same server code, with or without the TUI!
