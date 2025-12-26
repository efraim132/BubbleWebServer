package WebServer

import (
	"github.com/charmbracelet/lipgloss"
)

var serverStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

func Run(messages *[]string) string {
	//err := append(*messages,
	//	serverStyle.Render("Server: ")+"Hello World!",
	//)

	//if err != nil {
	//	// Figure out an error to put here
	//}
	return serverStyle.Render("Server: ") + "Hello World!"
}
