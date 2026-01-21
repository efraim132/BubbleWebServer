# BubbleServer Package

A simple, importable Go package that combines a BubbleTea TUI with an HTTP web server.

## Features

- Interactive TUI for server management
- Configurable commands
- HTTP route registration
- Built-in server control commands (start, stop, restart)
- Port configuration
- Request logging in TUI

## Installation

```bash

go get github.com/efraim132/BubbleWebServer/bubbleserver
```

## Quick Start

```go
package main

import (
    "github.com/efraim132/BubbleWebServer/bubbleserver"
    "fmt"
    "log"
    "net/http"
)

func main() {
    // Create a new server
    server := bubbleserver.NewServer(bubbleserver.Config{
        Banner:      "My Awesome Server",
        WelcomeText: "Welcome! Type /help for commands.",
        Debug:       false,
    })

    // Set default port
    server.SetDefaultPort(8080)

    // Register a custom TUI command
    server.RegisterCommand(bubbleserver.Command{
        Name:        "hello",
        Description: "Say hello",
        Handler: func(arg string) string {
            return "Hello from TUI!"
        },
    })

    // Register an HTTP route
    server.RegisterHTTPRoute(bubbleserver.HTTPRoute{
        Path: "/api/hello",
        Handler: func(w interface{}, req interface{}) {
            writer := w.(http.ResponseWriter)
            fmt.Fprintf(writer, "Hello from HTTP!\n")
        },
    })

    // Run the server
    if err := server.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### Creating a Server

```go
server := bubbleserver.NewServer(bubbleserver.Config{
    Banner:      string,  // ASCII art banner (optional)
    WelcomeText: string,  // Welcome message (optional)
    Debug:       bool,    // Show debug info (optional)
    LogFile:     string,  // Log file path (optional, defaults to "debug.log")
})
```

### Setting Default Port

```go
server.SetDefaultPort(8080)
```

### Registering TUI Commands

```go
server.RegisterCommand(bubbleserver.Command{
    Name:        "command-name",
    Description: "What this command does",
    Handler: func(arg string) string {
        // Command logic here
        return "Response message"
    },
})
```

Commands are invoked in the TUI by typing `/command-name`.

### Registering HTTP Routes

```go
server.RegisterHTTPRoute(bubbleserver.HTTPRoute{
    Path: "/your/path",
    Handler: func(w interface{}, req interface{}) {
        writer := w.(http.ResponseWriter)
        request := req.(*http.Request)

        // Handle the HTTP request
        fmt.Fprintf(writer, "Response")
    },
})
```

### Running the Server

```go
if err := server.Run(); err != nil {
    log.Fatal(err)
}
```

## Built-in Commands

The following commands are available by default:

- `/start server` - Start the HTTP server
- `/stop server` - Stop the HTTP server
- `/restart server` - Restart the HTTP server
- `/set port` - Change the server port
- `/help` - Show all available commands
- `done` - Exit the application

## TUI Controls

- Type commands in the input area at the bottom
- Press Enter to execute
- Type `done` to quit
- Press Ctrl+C or Esc twice to force quit
- Scroll through messages with arrow keys

## Example

See `example_main.go` in the repository root for a complete example.
