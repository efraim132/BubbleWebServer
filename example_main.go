//go:build example
// +build example

package examples

import (
	"fmt"
	"github.com/efraim132/BubbleWebServer/bubbleserver"
	"log"
	"net/http"
)

func RunExample() {
	// Create a new server with custom configuration
	server := bubbleserver.NewServer(bubbleserver.Config{
		Banner:      " /$$      /$$           /$$       /$$$$$$$            /$$       /$$       /$$                    \n| $$  /$ | $$          | $$      | $$__  $$          | $$      | $$      | $$                    \n| $$ /$$$| $$  /$$$$$$ | $$$$$$$ | $$  \\ $$ /$$   /$$| $$$$$$$ | $$$$$$$ | $$  /$$$$$$   /$$$$$$$\n| $$/$$ $$ $$ /$$__  $$| $$__  $$| $$$$$$$ | $$  | $$| $$__  $$| $$__  $$| $$ /$$__  $$ /$$_____/\n| $$$$_  $$$$| $$$$$$$$| $$  \\ $$| $$__  $$| $$  | $$| $$  \\ $$| $$  \\ $$| $$| $$$$$$$$|  $$$$$$ \n| $$$/ \\  $$$| $$_____/| $$  | $$| $$  \\ $$| $$  | $$| $$  | $$| $$  | $$| $$| $$_____/ \\____  $$\n| $$/   \\  $$|  $$$$$$$| $$$$$$$/| $$$$$$$/|  $$$$$$/| $$$$$$$/| $$$$$$$/| $$|  $$$$$$$ /$$$$$$$/\n|__/     \\__/ \\_______/|_______/ |_______/  \\______/ |_______/ |_______/ |__/ \\_______/|_______/ \n",
		WelcomeText: "Written by Efraim v0.1C\n\nWelcome to the control panel!\nType a message and press Enter to send.\nUse /help to see available commands.",
		Debug:       true,
	})

	// Set default port
	server.SetDefaultPort(8080)

	// Register custom TUI commands
	server.RegisterCommand(bubbleserver.Command{
		Name:        "greet",
		Description: "Send a greeting message",
		Handler: func(arg string) string {
			return "Hello from custom command!"
		},
	})

	server.RegisterCommand(bubbleserver.Command{
		Name:        "status",
		Description: "Show server status",
		Handler: func(arg string) string {
			return "Server is running and ready!"
		},
	})

	// Register HTTP routes
	server.RegisterHTTPRoute(bubbleserver.HTTPRoute{
		Path: "/hello",
		Handler: func(w interface{}, req interface{}) {
			writer := w.(http.ResponseWriter)
			fmt.Fprintf(writer, "Hello, World!\n")
		},
	})

	server.RegisterHTTPRoute(bubbleserver.HTTPRoute{
		Path: "/headers",
		Handler: func(w interface{}, req interface{}) {
			writer := w.(http.ResponseWriter)
			request := req.(*http.Request)
			for name, headers := range request.Header {
				for _, h := range headers {
					fmt.Fprintf(writer, "%v: %v\n", name, h)
				}
			}
		},
	})

	server.RegisterHTTPRoute(bubbleserver.HTTPRoute{
		Path: "/api/data",
		Handler: func(w interface{}, req interface{}) {
			writer := w.(http.ResponseWriter)
			writer.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(writer, `{"message": "This is a custom API endpoint", "status": "ok"}`)
		},
	})

	// Run the server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
