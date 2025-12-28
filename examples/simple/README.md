# Simple Web Server Example

This example demonstrates a basic web server using BubbleWebServer with different response types (HTML, JSON, Text) and automatic error handling.

## What This Example Shows

- HTML responses (web pages)
- JSON responses (API endpoints)
- Plain text responses
- Error handling
- Request logging with callbacks
- Lifecycle hooks (OnStart, OnStop)
- Graceful shutdown

## Running the Example

```bash
cd examples/simple
go run main.go
```

The server will start on **http://localhost:8080**

<!-- 📸 SCREENSHOT: Browser showing the home page or terminal showing curl responses -->
<!-- Show either the browser viewing localhost:8080 or terminal with curl outputs -->
<!-- Suggested filename: ../../assets/simple-example-browser.png -->
![Simple Example in Browser](../../assets/simple-example-browser.png)

## Endpoints

### 1. Home Page - `GET /`
Returns an HTML web page with links to all endpoints.

**Try it:**
```bash
# In your browser
http://localhost:8080/

# Or with curl
curl http://localhost:8080/
```

### 2. Hello - `GET /hello`
Returns a simple plain text response.

**Try it:**
```bash
curl http://localhost:8080/hello
```

**Response:**
```
Hello from BubbleWebServer!
```

### 3. Users API - `GET /api/users`
Returns JSON data with a list of users.

**Try it:**
```bash
curl http://localhost:8080/api/users
```

**Response:**
```json
{
  "success": true,
  "count": 3,
  "users": [
    {"id": 1, "name": "Alice", "email": "alice@example.com"},
    {"id": 2, "name": "Bob", "email": "bob@example.com"},
    {"id": 3, "name": "Charlie", "email": "charlie@example.com"}
  ]
}
```

### 4. Error Example - `GET /api/error`
Demonstrates automatic error handling.

**Try it:**
```bash
curl http://localhost:8080/api/error
```

**Response:**
```json
{
  "error": "The requested resource was not found",
  "code": "NOT_FOUND"
}
```

**Status Code:** 404 Not Found

## Code Walkthrough

### Server Setup

```go
srv := server.New().                                    // Create new server
    Port(8080).                                         // Set port
    Route("/", homeHandler).                            // Add routes
    Route("/hello", helloHandler).
    Route("/api/users", usersHandler).
    Route("/api/error", errorHandler).
    OnRequest(func(info server.RequestInfo) {          // Log each request
        log.Printf("[%s] %s from %s", info.Method, info.Path, info.IP)
    }).
    OnStart(func() {                                    // Called when server starts
        log.Println("Server started successfully!")
    }).
    OnStop(func() {                                     // Called when server stops
        log.Println("Server stopped gracefully")
    }).
    Build()
```

### HTML Response

```go
func homeHandler(w http.ResponseWriter, r *http.Request) error {
    html := `<!DOCTYPE html>...<h1>Hello</h1>...`
    return server.HTML(w, 200, html)
}
```

### Text Response

```go
func helloHandler(w http.ResponseWriter, r *http.Request) error {
    return server.Text(w, 200, "Hello from BubbleWebServer!\n")
}
```

### JSON Response

```go
func usersHandler(w http.ResponseWriter, r *http.Request) error {
    users := []map[string]interface{}{
        {"id": 1, "name": "Alice", "email": "alice@example.com"},
        // ...
    }
    return server.JSON(w, 200, map[string]interface{}{
        "success": true,
        "count":   len(users),
        "users":   users,
    })
}
```

### Error Handling

```go
func errorHandler(w http.ResponseWriter, r *http.Request) error {
    // Just return an error - the server handles the rest!
    return server.NotFoundError("The requested resource was not found")
}
```

## Key Concepts Demonstrated

### 1. Response Helpers
Instead of manually setting headers and status codes, use helpers:
- `server.HTML(w, 200, html)` - HTML response
- `server.JSON(w, 200, data)` - JSON response
- `server.Text(w, 200, text)` - Text response

### 2. Automatic Error Handling
Return errors from handlers:
```go
return server.NotFoundError("message")  // Becomes 404 with JSON error
```

No need to manually write error responses!

### 3. Request Callbacks
Use `OnRequest()` to log or monitor every request:
```go
OnRequest(func(info server.RequestInfo) {
    log.Printf("[%s] %s from %s", info.Method, info.Path, info.IP)
})
```

### 4. Lifecycle Hooks
Run code when the server starts or stops:
```go
OnStart(func() {
    log.Println("Server started!")
})
OnStop(func() {
    log.Println("Server stopped!")
})
```

### 5. Graceful Shutdown
The server waits for Ctrl+C, then:
1. Receives the signal
2. Stops accepting new requests
3. Waits for active requests to finish
4. Shuts down cleanly

## Customizing This Example

### Add a New Endpoint

1. Add the route:
```go
Route("/goodbye", goodbyeHandler).
```

2. Create the handler:
```go
func goodbyeHandler(w http.ResponseWriter, r *http.Request) error {
    return server.Text(w, 200, "Goodbye!\n")
}
```

### Change the Port

```go
Port(3000).  // Instead of 8080
```

### Add Custom Logic

```go
OnRequest(func(info server.RequestInfo) {
    // Your custom logging/monitoring logic
    if strings.HasPrefix(info.Path, "/api/") {
        log.Println("API request:", info.Path)
    }
})
```

## Next Steps

- Check out **`examples/api/`** for a RESTful API example with CRUD operations
- Check out **`examples/tui/`** for using the Terminal UI
- Read the main **README.md** for complete documentation

## Questions?

This example is heavily commented - read through `main.go` for detailed explanations of how everything works!
