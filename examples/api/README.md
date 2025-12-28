# RESTful API Example - Todo API

This example demonstrates a complete RESTful API using BubbleWebServer. It implements CRUD operations (Create, Read, Update, Delete) for a Todo list application.

## What This Example Shows

- RESTful API design patterns
- CRUD operations
- JSON request/response handling
- URL parameter extraction
- HTTP method validation
- Input validation
- Custom error messages
- In-memory data storage

## Running the Example

```bash
cd examples/api
go run main.go
```

The server will start on **http://localhost:3000**

<!-- 📸 SCREENSHOT: API responses in Postman/Insomnia or curl -->
<!-- Show curl commands and their JSON responses, or a REST client interface -->
<!-- Suggested filename: ../../assets/api-example-responses.png -->
![API Example Responses](../../assets/api-example-responses.png)

## API Endpoints

### 1. Health Check - `GET /health`
Check if the API is running.

**Request:**
```bash
curl http://localhost:3000/health
```

**Response:**
```json
{
  "status": "healthy"
}
```

### 2. List All Todos - `GET /api/todos`
Get all todos.

**Request:**
```bash
curl http://localhost:3000/api/todos
```

**Response:**
```json
{
  "todos": [
    {"id": 1, "title": "Learn Go", "completed": true},
    {"id": 2, "title": "Build a web server", "completed": true},
    {"id": 3, "title": "Deploy to production", "completed": false}
  ],
  "count": 3
}
```

### 3. Get Single Todo - `GET /api/todos/:id`
Get a specific todo by ID.

**Request:**
```bash
curl http://localhost:3000/api/todos/1
```

**Response:**
```json
{
  "id": 1,
  "title": "Learn Go",
  "completed": true
}
```

**Error (404):**
```bash
curl http://localhost:3000/api/todos/999
```
```json
{
  "error": "Todo with ID 999 not found",
  "code": "NOT_FOUND"
}
```

### 4. Create Todo - `POST /api/todos/create`
Create a new todo.

**Request:**
```bash
curl -X POST http://localhost:3000/api/todos/create \
  -H "Content-Type: application/json" \
  -d '{"title":"Write documentation"}'
```

**Response (201 Created):**
```json
{
  "id": 4,
  "title": "Write documentation",
  "completed": false
}
```

**Validation Error (400):**
```bash
curl -X POST http://localhost:3000/api/todos/create \
  -H "Content-Type: application/json" \
  -d '{"title":""}'
```
```json
{
  "error": "Title is required",
  "code": "BAD_REQUEST"
}
```

### 5. Delete Todo - `DELETE /api/todos/:id`
Delete a todo by ID.

**Request:**
```bash
curl -X DELETE http://localhost:3000/api/todos/1
```

**Response:**
```json
{
  "message": "Todo 1 deleted"
}
```

## RESTful Design Patterns

### Resource-Based URLs
```
/api/todos       - Collection of todos
/api/todos/:id   - Single todo resource
```

### HTTP Methods Map to Actions
```
GET    /api/todos       - List all (Read)
POST   /api/todos/create - Create new (Create)
GET    /api/todos/:id   - Get one (Read)
DELETE /api/todos/:id   - Delete one (Delete)
```

### Status Codes
```
200 OK              - Successful GET/DELETE
201 Created         - Successful POST (resource created)
400 Bad Request     - Invalid input
404 Not Found       - Resource doesn't exist
405 Method Not Allowed - Wrong HTTP method
```

## Code Walkthrough

### Data Model

```go
type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
```

The `json:"..."` tags tell Go how to convert struct fields to JSON.

### Server Setup

```go
srv := server.New().
    Port(3000).
    Route("/api/todos", listTodosHandler).
    Route("/api/todos/", todoHandler).  // Handles /:id
    Route("/api/todos/create", createTodoHandler).
    Route("/health", healthHandler).
    Build()
```

### List Todos Handler

```go
func listTodosHandler(w http.ResponseWriter, r *http.Request) error {
    // Validate HTTP method
    if r.Method != http.MethodGet {
        return server.ErrMethodNotAllowed
    }

    // Return data
    return server.JSON(w, 200, map[string]interface{}{
        "todos": todos,
        "count": len(todos),
    })
}
```

### Create Todo Handler

```go
func createTodoHandler(w http.ResponseWriter, r *http.Request) error {
    // Validate method
    if r.Method != http.MethodPost {
        return server.ErrMethodNotAllowed
    }

    // Parse JSON request
    var req struct {
        Title string `json:"title"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return server.BadRequestError("Invalid JSON")
    }

    // Validate input
    if req.Title == "" {
        return server.BadRequestError("Title is required")
    }

    // Create todo
    todo := Todo{
        ID:        nextID,
        Title:     req.Title,
        Completed: false,
    }
    nextID++
    todos = append(todos, todo)

    // Return 201 Created
    return server.Created(w, todo)
}
```

### URL Parameter Extraction

```go
func todoHandler(w http.ResponseWriter, r *http.Request) error {
    // Extract ID from URL: /api/todos/123 -> "123"
    path := strings.TrimPrefix(r.URL.Path, "/api/todos/")

    // Convert to integer
    id, err := strconv.Atoi(path)
    if err != nil {
        return server.BadRequestError("Invalid todo ID")
    }

    // Find the todo
    var todo Todo
    for _, t := range todos {
        if t.ID == id {
            todo = t
            break
        }
    }

    // Handle different methods
    switch r.Method {
    case http.MethodGet:
        return server.JSON(w, 200, todo)
    case http.MethodDelete:
        // Delete logic...
        return server.JSON(w, 200, map[string]string{
            "message": "Todo deleted",
        })
    default:
        return server.ErrMethodNotAllowed
    }
}
```

## Key Concepts Demonstrated

### 1. JSON Request Parsing

```go
var req struct {
    Title string `json:"title"`
}
json.NewDecoder(r.Body).Decode(&req)
```

### 2. Input Validation

```go
if req.Title == "" {
    return server.BadRequestError("Title is required")
}
```

### 3. HTTP Method Validation

```go
if r.Method != http.MethodPost {
    return server.ErrMethodNotAllowed
}
```

### 4. URL Parameters

```go
path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
id, err := strconv.Atoi(path)
```

### 5. Appropriate Status Codes

```go
return server.Created(w, todo)   // 201
return server.JSON(w, 200, data) // 200
return server.ErrNotFound        // 404
```

## Testing the API

### Using curl

```bash
# List todos
curl http://localhost:3000/api/todos

# Create todo
curl -X POST http://localhost:3000/api/todos/create \
  -H "Content-Type: application/json" \
  -d '{"title":"Test todo"}'

# Get specific todo
curl http://localhost:3000/api/todos/1

# Delete todo
curl -X DELETE http://localhost:3000/api/todos/1
```

### Using HTTP Clients

- **Postman** - Import endpoints and test
- **Insomnia** - REST client
- **HTTPie** - `http POST localhost:3000/api/todos/create title="Test"`

## Extending This Example

### Add UPDATE Operation

```go
// Add route
Route("/api/todos/update", updateTodoHandler).

// Handler
func updateTodoHandler(w http.ResponseWriter, r *http.Request) error {
    if r.Method != http.MethodPut {
        return server.ErrMethodNotAllowed
    }

    var req struct {
        ID        int    `json:"id"`
        Title     string `json:"title"`
        Completed bool   `json:"completed"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return server.BadRequestError("Invalid JSON")
    }

    // Find and update todo...
    return server.JSON(w, 200, updatedTodo)
}
```

### Add Database

Replace the in-memory `todos` slice with a real database:
```go
// Instead of: todos = append(todos, todo)
db.Create(&todo)

// Instead of: for _, t := range todos
todos, err := db.Find()
```

### Add Authentication

```go
func requireAuth(handler server.HandlerFunc) server.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) error {
        token := r.Header.Get("Authorization")
        if token == "" {
            return server.ErrUnauthorized
        }
        return handler(w, r)
    }
}

// Use it
Route("/api/todos/create", requireAuth(createTodoHandler)).
```

## Production Considerations

### Current State (Example)
- In-memory storage (data lost on restart)
- No authentication
- No pagination
- No rate limiting

### Production Ready
- Database (PostgreSQL, MySQL, MongoDB)
- Authentication & authorization
- Input sanitization
- Pagination for large lists
- Rate limiting
- Logging & monitoring
- CORS headers for web clients

## Next Steps

- Check out **`examples/simple/`** for basic web server patterns
- Check out **`examples/tui/`** for using the Terminal UI
- Read the main **README.md** for complete documentation
- Explore Phase 2 features for middleware and advanced patterns

## Questions?

This example is heavily commented - read through `main.go` for detailed explanations of every REST pattern used!
