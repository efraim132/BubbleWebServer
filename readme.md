# BubbleWebServer

A production-ready, developer-friendly HTTP server package for Go with an optional Terminal User Interface (TUI). Built to be your baseline web server for all future projects - just import, define endpoints, and go!

```
 /$$      /$$           /$$       /$$$$$$$            /$$       /$$       /$$
| $$  /$ | $$          | $$      | $$__  $$          | $$      | $$      | $$
| $$ /$$$| $$  /$$$$$$ | $$$$$$$ | $$  \ $$ /$$   /$$| $$$$$$$ | $$$$$$$ | $$  /$$$$$$   /$$$$$$$
| $$/$$ $$ $$ /$$__  $$| $$__  $$| $$$$$$$ | $$  | $$| $$__  $$| $$__  $$| $$ /$$__  $$ /$$_____/
| $$$$_  $$$$| $$$$$$$$| $$  \ $$| $$__  $$| $$  | $$| $$  \ $$| $$  \ $$| $$| $$$$$$$$|  $$$$$$
| $$$/ \  $$$| $$_____/| $$  | $$| $$  \ $$| $$  | $$| $$  | $$| $$  | $$| $$| $$_____/ \____  $$
| $$/   \  $$|  $$$$$$$| $$$$$$$/| $$$$$$$/|  $$$$$$/| $$$$$$$/| $$$$$$$/| $$|  $$$$$$$ /$$$$$$$/
|__/     \__/ \_______/|_______/ |_______/  \______/ |_______/ |_______/ |__/ \_______/|_______/
```

## Features

FeaturesPhase 1 - Core Functionality (Complete)

- Fluent API - Chain methods for clean, readable configuration
- Response Helpers - One-line JSON, HTML, Text, and Error responses
- Error Handling - Return errors from handlers, automatic conversion to HTTP responses
- Importable Package - Use in any project, zero dependencies on TUI
- Terminal UI (Optional) - Beautiful TUI for development with real-time request monitoring
- Comprehensive Tests - Full test coverage for reliability
- Detailed Documentation - Every file heavily commented for learning

### What Makes This Different

Unlike other Go web frameworks, BubbleWebServer is designed specifically for:
- **Rapid prototyping** - Get a server running in 3 lines of code
- **Learning** - Every file is heavily documented with explanations
- **Reusability** - Import into any project without modification
- **Development UX** - Optional TUI for visual feedback during development
- **Simplicity** - No magic, no complex routing, just clean Go code

## Quick Start

### Installation

```bash
# Clone or download this repository
git clone <your-repo-url>
cd BubbleWebServer

# Install dependencies
go mod download
```

### Simple Server (3 lines!)

```go
import "BubbleWebServer/server"

func main() {
    server.New().
        Route("/hello", func(w http.ResponseWriter, r *http.Request) error {
            return server.Text(w, 200, "Hello, World!")
        }).
        Build().
        Start()
}
```

### With TUI (Development Mode)

```go
func main() {
    srv := server.New().
        Port(8080).
        Route("/hello", helloHandler).
        Build()

    srv.StartWithTUI()  // Launch with Terminal UI
}
```

<!-- 📸 SCREENSHOT 1: TUI Main Interface -->
<!-- Show the TUI running with the banner, command prompt, and maybe one or two request logs -->
<!-- Suggested filename: assets/tui-main-screen.png -->
![TUI Main Screen](assets/tui-main-screen.png)

<!-- 📸 SCREENSHOT 2: Port Selection Form -->
<!-- Show the centered port selection form with the rounded border -->
<!-- Suggested filename: assets/tui-port-selection.png -->
![Port Selection Form](assets/tui-port-selection.png)

## Core Concepts

### 1. Fluent API Pattern

Chain methods together for readable configuration:

```go
server.New().
    Port(8080).
    Host("localhost").
    Route("/users", getUsersHandler).
    Route("/products", getProductsHandler).
    OnRequest(logRequest).
    Build()
```

### 2. Response Helpers

Stop writing boilerplate - use helpers:

```go
// Before
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(200)
json.NewEncoder(w).Encode(data)

// After
return server.JSON(w, 200, data)
```

Available helpers:
- `JSON(w, status, data)` - JSON response
- `Text(w, status, text)` - Plain text
- `HTML(w, status, html)` - HTML page
- `Error(w, status, message)` - Error response
- `Created(w, data)` - 201 Created
- `OK(w, data)` - 200 OK
- `BadRequest(w, msg)` - 400 Bad Request
- `NotFound(w, msg)` - 404 Not Found
- And more...

### 3. Error Handling

Return errors from handlers - they're automatically converted:

```go
func getUser(w http.ResponseWriter, r *http.Request) error {
    user := findUser(id)
    if user == nil {
        return server.ErrNotFound  // Becomes 404 JSON response
    }
    return server.JSON(w, 200, user)
}
```

Predefined errors:
- `ErrBadRequest` (400)
- `ErrUnauthorized` (401)
- `ErrForbidden` (403)
- `ErrNotFound` (404)
- `ErrInternalServer` (500)

Custom errors:
```go
return server.NotFoundError("User with ID 123 not found")
return server.BadRequestError("Email is required")
```

### 4. Terminal UI

Get visual feedback during development:

```go
srv.StartWithTUI()  // Instead of srv.Start()
```

TUI Commands:
- `/start` - Start the server
- `/stop` - Stop the server
- `/restart` - Restart the server
- `/port` - Change server port (interactive form)
- `/status` - Show server status
- `/help` - Show help
- `quit` - Exit

## Examples

### Basic Web Server
```go
// examples/simple/main.go
srv := server.New().
    Port(8080).
    Route("/", homeHandler).
    Route("/api/hello", apiHandler).
    Build()

srv.StartWithTUI()
```

### RESTful API
```go
// examples/api/main.go
srv := server.New().
    Port(3000).
    Route("/api/todos", listTodosHandler).
    Route("/api/todos/create", createTodoHandler).
    Route("/api/todos/", todoHandler).
    Build()

srv.Start()  // Headless mode for production
```

See the `examples/` directory for complete, runnable examples:
- **`examples/simple/`** - Basic web server with HTML, text, and JSON
- **`examples/api/`** - RESTful Todo API with CRUD operations
- **`examples/tui/`** - Using the TUI in a new project

## Project Structure

```
BubbleWebServer/
├── server/              # Core server package (importable!)
│   ├── server.go       # Main server with fluent API
│   ├── response.go     # Response helper functions
│   ├── errors.go       # Error handling system
│   ├── tui.go          # Terminal UI integration
│   ├── *_test.go       # Comprehensive tests
│   └── ...
├── examples/           # Example applications
│   ├── simple/         # Basic web server
│   ├── api/            # RESTful API
│   └── tui/            # TUI example
├── main.go             # TUI demo application
├── go.mod              # Go module definition
└── README.md           # This file
```

## Usage in New Projects

### Step 1: Import the Package

```go
import "BubbleWebServer/server"
```

### Step 2: Create Your Server

```go
srv := server.New().
    Port(8080).
    Route("/endpoint", yourHandler).
    Build()
```

### Step 3: Start It

```go
// With TUI (development)
srv.StartWithTUI()

// Or without TUI (production)
srv.Start()
```

That's it! You have a fully functional web server.

## Configuration Options

### TUI Customization

```go
config := server.DefaultTUIConfig()
config.Banner = "My Custom Banner"
config.Title = "My Server Control Panel"
config.ShowDebug = true
config.BannerColor = "5"  // Lipgloss color
config.ServerColor = "2"
config.SystemColor = "1"

srv.StartWithTUIConfig(config)
```

### Server Options

```go
srv := server.New().
    Port(8080).              // Set port
    Host("localhost").       // Set host (default: 0.0.0.0)
    Route("/path", handler). // Add route
    OnRequest(callback).     // Request callback
    OnStart(callback).       // Start callback
    OnStop(callback).        // Stop callback
    Build()
```

## API Reference

### Server Methods

- `New()` - Create new server
- `Port(int)` - Set port (fluent)
- `Host(string)` - Set host (fluent)
- `Route(path, handler)` - Add route (fluent)
- `OnRequest(callback)` - Set request callback (fluent)
- `OnStart(callback)` - Set start callback (fluent)
- `OnStop(callback)` - Set stop callback (fluent)
- `Build()` - Build server (fluent)
- `Start()` - Start server
- `StartWithTUI()` - Start with TUI
- `StartWithTUIConfig(config)` - Start with custom TUI
- `Stop()` - Stop server
- `Restart()` - Restart server
- `IsRunning()` - Check if running
- `GetPort()` - Get port
- `GetHost()` - Get host

### Response Helpers

See `server/response.go` for full list of helpers.

### Error Types

See `server/errors.go` for full list of error types and functions.

## Testing

Run all tests:

```bash
go test ./server/...
```

Run with coverage:

```bash
go test ./server/... -cover
```

Run with verbose output:

```bash
go test ./server/... -v
```

## Dependencies

- **Runtime**: Go standard library only (no external dependencies for headless mode)
- **TUI Mode**: charmbracelet/bubbletea, bubbles, lipgloss, huh (only if using TUI)

## Development

### Running the Demo TUI

```bash
go run main.go
```

Then type `/start` to start the server and visit http://localhost:8090/hello

### Running Examples

```bash
cd examples/simple
go run main.go
```


## Contributing

This is a personal project designed for learning and rapid prototyping. Feel free to fork and adapt for your needs!

## License

This project is provided as-is for educational and personal use.

## Author

Created by Efraim - Version 0.1A

