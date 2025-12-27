package WebServer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var port = 8090

// HTTPRequestMsg is sent to the BubbleTea UI when an HTTP request is received
type HTTPRequestMsg struct {
	Path   string
	Method string
}

type webCommand struct {
	webPath string
	method  func(http.ResponseWriter, *http.Request)
}

var isWebServerRunning = false

var serverStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

var program *tea.Program
var server *http.Server

var webCommands []webCommand = []webCommand{}

//Utility ==============================================================================================================

// SetProgram sets the BubbleTea program instance so we can send messages to it
func SetProgram(p *tea.Program) {
	program = p
}

func stylizeServerMessage(message string) string {
	return serverStyle.Render("Server: ") + message
}

// webCommands handling ================================================================================================

func addCommand(cmd webCommand) error {
	webCommands = append(webCommands, cmd)
	return nil
}

// Webserver handling ==================================================================================================

func WebInit() error {
	err := addCommand(webCommand{
		webPath: "/hello",
		method:  hello,
	})
	if err != nil {
		return fmt.Errorf("web server startup failed: %v", err)
	}
	err = addCommand(webCommand{
		webPath: "/header",
		method:  headers,
	})
	if err != nil {
		return fmt.Errorf("web server startup failed: %v", err)
	}

	return nil

}

func StartWebserverGoRoutine(rPort int) string {
	if isWebServerRunning {
		return stylizeServerMessage("Webserver is already running!")
	}
	if rPort != -1 {
		port = rPort
	}
	go main()
	isWebServerRunning = true
	return stylizeServerMessage(fmt.Sprintf("Started Up Webserver! Listening to port %d", port))
}

func StopWebserver() string {
	if !isWebServerRunning {
		return stylizeServerMessage("Webserver is not running!")
	}

	if server != nil {
		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return stylizeServerMessage(fmt.Sprintf("Error stopping server: %v", err))
		}

		isWebServerRunning = false
		server = nil
		return stylizeServerMessage("Stopped Webserver!")
	}

	return stylizeServerMessage("Server reference not found!")
}

func RestartWebserver(rPort int) string {
	if StopWebserver() == "Webserver is not running!" {
		return stylizeServerMessage("Webserver is already stopped!")
	}
	StartWebserverGoRoutine(rPort) // Use existing port
	return stylizeServerMessage(fmt.Sprintf("Restarted Webserver on port %d!", port))
}

func hello(w http.ResponseWriter, req *http.Request) {
	// Send event to BubbleTea UI
	if program != nil {
		program.Send(HTTPRequestMsg{
			Path:   req.URL.Path,
			Method: req.Method,
		})
	}

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	// Send event to BubbleTea UI
	if program != nil {
		program.Send(HTTPRequestMsg{
			Path:   req.URL.Path,
			Method: req.Method,
		})
	}

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	mux := http.NewServeMux()

	for _, webCommand := range webCommands {
		mux.HandleFunc(webCommand.webPath, webCommand.method)
	}

	server = &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", port),
		Handler: mux,
	}

	// ListenAndServe blocks, so this runs in the goroutine
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Send error message to UI if server fails to start
		if program != nil {
			program.Send(HTTPRequestMsg{
				Path:   "ERROR",
				Method: fmt.Sprintf("Server error: %v", err),
			})
		}
	}
}
