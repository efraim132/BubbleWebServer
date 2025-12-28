# Terminal UI (TUI) Example

This example demonstrates how to use BubbleWebServer with the built-in Terminal User Interface in a new project. The TUI provides a beautiful, interactive terminal interface for controlling and monitoring your server during development.

## What This Example Shows

- Starting a server with TUI in one line
- Real-time request monitoring in the terminal
- Interactive server control commands
- Custom TUI configuration
- Server lifecycle management through TUI

## Running the Example

```bash
cd examples/tui
go run main.go
```

You'll see a terminal interface with:
- A banner
- Server status information
- A command prompt (`>`)

## Using the TUI

### Starting the Server

When the TUI launches, the HTTP server is **not running** yet. To start it:

1. Type: `/start`
2. Press Enter
3. You'll see: "Server: Started on port 8080"

### Testing the Server

Open another terminal and try the endpoints:

```bash
# Test the home page
curl http://localhost:8080/

# Test the API
curl http://localhost:8080/api/hello
```

**Watch the TUI** - you'll see each request appear in real-time!

```
Server: Received GET request to / from 127.0.0.1:54321
Server: Received GET request to /api/hello from 127.0.0.1:54322
```

### TUI Commands

Type these commands in the TUI prompt:

| Command | Description |
|---------|-------------|
| `/start` | Start the HTTP server |
| `/stop` | Stop the HTTP server |
| `/restart` | Restart the HTTP server |
| `/port` | Change server port (interactive form with validation) |
| `/status` | Show current server status |
| `/help` | Show help message |
| `quit` or `exit` | Exit the TUI (server stops automatically) |

## Code Walkthrough

### Option 1: Default TUI (Simple)

```go
srv := server.New().
    Port(8080).
    Route("/", homeHandler).
    Route("/api/hello", apiHelloHandler).
    Build()

// Start with default TUI
srv.StartWithTUI()
```

That's it! You get a TUI with default settings.

### Option 2: Custom TUI (Advanced)

```go
// Create custom configuration
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
config.ShowDebug = true  // Show debug info

srv := server.New().
    Port(8080).
    Route("/", homeHandler).
    Build()

// Start with custom TUI
srv.StartWithTUIConfig(config)
```

### TUI Configuration Options

```go
type TUIConfig struct {
    Banner      string // ASCII art banner
    Title       string // Title shown in TUI
    ShowDebug   bool   // Show debug information
    BannerColor string // Lipgloss color for banner
    ServerColor string // Lipgloss color for server messages
    SystemColor string // Lipgloss color for system messages
}
```

## What You'll See

### Initial Screen

```
 __  __         ____
|  \/  |_   _  / ___|  ___ _ ____   _____ _ __
| |\/| | | | | \___ \ / _ \ '__\ \ / / _ \ '__|
| |  | | |_| |  ___) |  __/ |   \ V /  __/ |
|_|  |_|\__, | |____/ \___|_|    \_/ \___|_|
        |___/

My Custom Server

Server ready on port 8080
Commands: /start, /stop, /restart, /help, quit

> _
```

<!-- 📸 SCREENSHOT: Initial TUI Screen -->
<!-- Show the TUI when first launched, with banner and command prompt -->
<!-- Suggested filename: ../../assets/tui-initial-screen.png -->
![TUI Initial Screen](../../assets/tui-initial-screen.png)

### After Starting Server

```
> /start
Server: Started on port 8080

> _
```

### When Requests Come In

```
> /start
Server: Started on port 8080
Server: Received GET request to / from 127.0.0.1:54321
Server: Received GET request to /api/hello from 127.0.0.1:54322

> _
```

<!-- 📸 SCREENSHOT: TUI with Request Logs -->
<!-- Show the TUI with several request logs displayed, showing real-time monitoring -->
<!-- Suggested filename: ../../assets/tui-request-logs.png -->
![TUI Request Monitoring](../../assets/tui-request-logs.png)

## Changing the Port

The TUI includes an interactive port selection feature with validation:

### Using Port Selection

1. Type: `/port` in the TUI prompt
2. Press Enter
3. You'll see a centered form asking for a port number
4. Enter a port (e.g., `3000`)
5. Press Enter to confirm

### Port Validation

The form validates your input:
- **Must be a number** - Letters and symbols are rejected
- **Range check** - Must be between 1 and 65535
- **Note** - Ports below 1024 may require admin/root privileges on your system

### Example

```
> /port
(Form appears centered on screen)

┌─ Set Server Port ──┐
│                    │
│  HTTP Port         │
│  > 3000            │
│                    │
└────────────────────┘

(After pressing Enter)
System: Port set to 3000
```

<!-- 📸 SCREENSHOT: Port Selection Form -->
<!-- Show the centered port selection form with validation -->
<!-- Suggested filename: ../../assets/tui-port-form.png -->
![Port Selection Form](../../assets/tui-port-form.png)

### Important Notes

- **Cannot change port while server is running** - You must `/stop` the server first
- **Port persists** - The new port is used when you `/start` the server
- **ESC to cancel** - Press ESC in the form to cancel without changing
- **Privileged ports** - Ports below 1024 (like 80, 443) may require admin/root privileges to bind

## When to Use TUI vs Regular Start

### Use `StartWithTUI()` for:

- **Development** - Visual feedback while coding
- **Debugging** - See exactly what's happening
- **Demos** - Looks professional in presentations
- **Local Testing** - Easy server control
- **Learning** - Understand server behavior visually

### Use `Start()` for:

- **Production** - Headless servers, cloud deployments
- **Docker Containers** - No terminal interaction needed
- **Background Services** - Runs silently
- **CI/CD Pipelines** - Automated deployments
- **Systemd Services** - Linux service management

## Benefits of TUI

### 1. Visual Feedback
See what's happening without digging through logs:
```
Server: Received GET request to /api/users from 127.0.0.1:54321
Server: Received POST request to /api/users/create from 127.0.0.1:54322
```

### 2. Easy Control
Start/stop with simple commands instead of Ctrl+C:
```
> /stop
Server: Stopped
> /start
Server: Started on port 8080
```

### 3. Better Debugging
See request patterns in real-time:
```
Server: Received GET request to /api/users
Server: Received GET request to /api/users/1
Server: Received GET request to /api/users/2
```

### 4. Professional Look
Impress colleagues with a clean terminal interface!

## Customizing the TUI

### Custom Banner

Create ASCII art at [patorjk.com/software/taag](https://patorjk.com/software/taag/):

```go
config.Banner = `
  ___      _                      ___    ___  ___
 / _ \    | |                    / _ \  / _ \|_  |
/ /_\ \_ _| |_ ___  ___ _ __ ___/ /_\ \/ /_\ \ | |
|  _  | '| | __/ _ \/ __| '_ ` _ \  _  |  _  | | |
| | | | |_| || (_) \__ \ | | | | | | | | | | |_| |_
\_| |_/\__,_|\__\___/|___/_| |_| |_\_| |_\_| |_(_)
`
```

### Custom Colors

Lipgloss color options:
```go
config.BannerColor = "5"   // Magenta
config.ServerColor = "2"   // Green
config.SystemColor = "1"   // Red
```

Or use hex colors:
```go
config.BannerColor = "#FF5733"
```

### Show Debug Info

```go
config.ShowDebug = true
```

Shows:
```
Debug: 120x40 | Messages: 15
```

## Combining with Other Examples

You can add TUI to any server:

### From Simple Example
```go
// examples/simple with TUI
srv := server.New().
    Port(8080).
    Route("/", homeHandler).
    Route("/hello", helloHandler).
    Route("/api/users", usersHandler).
    Build()

srv.StartWithTUI()  // Add this line!
```

### From API Example
```go
// examples/api with TUI
srv := server.New().
    Port(3000).
    Route("/api/todos", listTodosHandler).
    Route("/api/todos/create", createTodoHandler).
    Route("/api/todos/", todoHandler).
    Build()

srv.StartWithTUI()  // Add this line!
```

## Switching Between TUI and Headless

### Development (with TUI)
```go
if os.Getenv("ENV") == "dev" {
    srv.StartWithTUI()
} else {
    srv.Start()
}
```

### Or with a flag
```go
useTUI := flag.Bool("tui", false, "Start with TUI")
flag.Parse()

if *useTUI {
    srv.StartWithTUI()
} else {
    srv.Start()
}
```

Run with: `go run main.go -tui`

## Troubleshooting

### TUI Not Showing
- Make sure you're running in a terminal, not as a background process
- Check your terminal supports ANSI colors

### Commands Not Working
- Commands must start with `/` (except `quit`)
- Type `/help` to see all commands

### Server Won't Start
- Check if port is already in use
- Try a different port: `.Port(8081)`

## Next Steps

- Try adding your own endpoints and watch them in the TUI
- Customize the banner with your project name
- Check out **`examples/simple/`** for basic patterns
- Check out **`examples/api/`** for RESTful API patterns
- Read the main **README.md** for complete documentation

## Questions?

The TUI code is in `server/tui.go` and is heavily commented - check it out to understand how it works under the hood!
