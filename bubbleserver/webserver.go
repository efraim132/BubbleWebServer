package bubbleserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HTTPRequestMsg is sent to the BubbleTea UI when an HTTP request is received
type HTTPRequestMsg struct {
	Path   string
	Method string
}

type webServerManager struct {
	isRunning   bool
	server      *http.Server
	program     *tea.Program
	routes      []HTTPRoute
	serverStyle lipgloss.Style
}

func newWebServerManager() *webServerManager {
	return &webServerManager{
		isRunning:   false,
		routes:      []HTTPRoute{},
		serverStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
}

func (w *webServerManager) setProgram(p *tea.Program) {
	w.program = p
}

func (w *webServerManager) init(routes []HTTPRoute) error {
	w.routes = routes
	return nil
}

func (w *webServerManager) start(port int) string {
	if w.isRunning {
		return w.stylizeMessage("Webserver is already running!")
	}

	go w.runServer(port)
	w.isRunning = true
	return w.stylizeMessage(fmt.Sprintf("Started Up Webserver! Listening to port %d", port))
}

func (w *webServerManager) stop() string {
	if !w.isRunning {
		return w.stylizeMessage("Webserver is not running!")
	}

	if w.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := w.server.Shutdown(ctx); err != nil {
			return w.stylizeMessage(fmt.Sprintf("Error stopping server: %v", err))
		}

		w.isRunning = false
		w.server = nil
		return w.stylizeMessage("Stopped Webserver!")
	}

	return w.stylizeMessage("Server reference not found!")
}

func (w *webServerManager) restart(port int) string {
	stopMsg := w.stop()
	if stopMsg == w.stylizeMessage("Webserver is not running!") {
		return w.stylizeMessage("Webserver is already stopped!")
	}
	return w.start(port)
}

func (w *webServerManager) runServer(port int) {
	mux := http.NewServeMux()

	// Register all routes
	for _, route := range w.routes {
		// Create a closure to capture route for each iteration
		currentRoute := route
		mux.HandleFunc(currentRoute.Path, func(wr http.ResponseWriter, req *http.Request) {
			// Send event to BubbleTea UI
			if w.program != nil {
				w.program.Send(HTTPRequestMsg{
					Path:   req.URL.Path,
					Method: req.Method,
				})
			}

			// Call the user's handler
			currentRoute.Handler(wr, req)
		})
	}

	w.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// ListenAndServe blocks
	if err := w.server.ListenAndServe(); err != http.ErrServerClosed {
		if w.program != nil {
			w.program.Send(HTTPRequestMsg{
				Path:   "ERROR",
				Method: fmt.Sprintf("Server error: %v", err),
			})
		}
	}
}

func (w *webServerManager) stylizeMessage(message string) string {
	return w.serverStyle.Render("Server: ") + message
}

// --- package-level default manager and compatibility wrappers ---
var defaultWebManager = newWebServerManager()

// WebInit registers default routes (/hello, /header) and initializes the manager
func WebInit() error {
	routes := []HTTPRoute{
		{
			Path: "/hello",
			Handler: func(w interface{}, req interface{}) {
				writer := w.(http.ResponseWriter)
				if defaultWebManager.program != nil {
					r := req.(*http.Request)
					defaultWebManager.program.Send(HTTPRequestMsg{Path: r.URL.Path, Method: r.Method})
				}
				fmt.Fprintf(writer, "hello\n")
			},
		},
		{
			Path: "/header",
			Handler: func(w interface{}, req interface{}) {
				writer := w.(http.ResponseWriter)
				if defaultWebManager.program != nil {
					r := req.(*http.Request)
					defaultWebManager.program.Send(HTTPRequestMsg{Path: r.URL.Path, Method: r.Method})
				}
				request := req.(*http.Request)
				for name, headers := range request.Header {
					for _, h := range headers {
						fmt.Fprintf(writer, "%v: %v\n", name, h)
					}
				}
			},
		},
	}
	return defaultWebManager.init(routes)
}

// SetProgram sets the BubbleTea program instance for the default manager
func SetProgram(p *tea.Program) {
	defaultWebManager.setProgram(p)
}

// StartWebserverGoRoutine starts the default webserver in a goroutine
func StartWebserverGoRoutine(rPort int) string {
	if rPort != -1 {
		return defaultWebManager.start(rPort)
	}
	return defaultWebManager.start(8090)
}

// StopWebserver stops the default webserver
func StopWebserver() string {
	return defaultWebManager.stop()
}

// RestartWebserver restarts the default webserver
func RestartWebserver(rPort int) string {
	return defaultWebManager.restart(rPort)
}
