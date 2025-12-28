# BubbleWebServer Examples

This directory contains examples showing how to use the BubbleWebServer package in different scenarios.

## Simple Example

A basic web server with multiple endpoints demonstrating different response types.

**Run:**
```bash
cd examples/simple
go run main.go
```

**Try:**
- http://localhost:8080/ - HTML page
- http://localhost:8080/hello - Text response
- http://localhost:8080/api/users - JSON response
- http://localhost:8080/api/error - Error handling

## API Example

A RESTful Todo API demonstrating CRUD operations with proper error handling.

**Run:**
```bash
cd examples/api
go run main.go
```

**Try:**
```bash
# List all todos
curl http://localhost:3000/api/todos

# Create a new todo
curl -X POST http://localhost:3000/api/todos/create \
  -H "Content-Type: application/json" \
  -d '{"title":"My new todo"}'

# Get a specific todo
curl http://localhost:3000/api/todos/1

# Delete a todo
curl -X DELETE http://localhost:3000/api/todos/1

# Health check
curl http://localhost:3000/health
```

## Key Features Demonstrated

### 1. Fluent API
```go
srv := server.New().
    Port(8080).
    Route("/hello", helloHandler).
    OnRequest(logHandler).
    Build()
```

### 2. Response Helpers
```go
// JSON response
server.JSON(w, 200, data)

// Text response
server.Text(w, 200, "Hello")

// HTML response
server.HTML(w, 200, "<h1>Hello</h1>")

// Error response
server.Error(w, 400, "Bad request")
```

### 3. Error Handling
```go
func handler(w http.ResponseWriter, r *http.Request) error {
    if somethingWrong {
        return server.ErrBadRequest  // Automatically handled!
    }
    return server.JSON(w, 200, data)
}
```

### 4. Lifecycle Hooks
```go
srv.OnStart(func() {
    log.Println("Server started!")
}).OnStop(func() {
    log.Println("Server stopped!")
})
```

## Creating Your Own Server

1. Import the package:
```go
import "BubbleWebServer/server"
```

2. Create and configure:
```go
srv := server.New().
    Port(8080).
    Route("/", homeHandler).
    Build()
```

3. Start it:
```go
srv.Start()
```

That's it! Your handlers can now use response helpers and return errors that are automatically handled.
