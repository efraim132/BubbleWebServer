package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

var banner = " __      __      ___.  __________     ___.  ___.   .__     \n" +
	"/  \\    /  \\ ____\\_ |__\\______   \\__ _\\_ |__\\_ |__ |  |   ____   ______\n" +
	"\\   \\/\\/   // __ \\| __ \\|    |  _/  |  \\ __ \\| __ \\|  | _/ __ \\ /  ___/\n" +
	" \\        /\\  ___/| \\_\\ \\    |   \\  |  / \\_\\ \\ \\_\\ \\  |_\\  ___/ \\___ \\ \n" +
	"  \\__/\\  /  \\___  >___  /______  /____/|___  /___  /____/\\___  >____  >\n" +
	"       \\/       \\/    \\/       \\/          \\/    \\/          \\/     \\/  \n"

func main() {
	fmt.Println(banner)
}

type Tcommand struct {
	function func(string) string
	cmd      string
}

type (
	errMsg         error
	screen         string
	stateChangeMsg screen
)
type server struct {
	port       int
	portString string
}

type Model struct {
	state       screen
	spinner     spinner.Model
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	statusStyle lipgloss.Style
	systemStyle lipgloss.Style
	portStyle   lipgloss.Style
	portForm    *huh.Form
	bannerStyle lipgloss.Style
	width       int
	height      int
	err         error
}

func initialModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Enter Command . . ."
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(banner + "Written by Efraim v0.1A" + gap + "Welcome to the control panel! \nType a message and press Enter to send.\n\"/start<or>stop server\"")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		state:       "chatting",
		spinner:     spinner.New(spinner.WithSpinner(spinner.Dot)),
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		statusStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		systemStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		portStyle: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")),
		portForm:    generateForm(),
		bannerStyle: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}),
		width:       30,
		height:      5,
		err:         nil,
	}
}

// Form stuff
func generateForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("HTTP Port").
				Placeholder("8080").
				Value(&portString). // Pointer to string variable
				// Validate function checks if it's a valid port number
				Validate(func(s string) error {
					// Check if it can be converted to integer
					port, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("port must be a number")
					}

					// Check if port is in valid range (1-65535)
					if port < 1 || port > 65535 {
						return fmt.Errorf("port must be between 1 and 65535")
					}

					// Optionally warn about privileged ports
					if port < 1024 {
						return fmt.Errorf("port %d requires admin privileges", port)
					}

					return nil // Validation passed
				}),
		),
	).WithWidth(60)

	return form
}
