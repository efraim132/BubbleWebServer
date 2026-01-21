package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	web "github.com/efrai/BubbleWebServer/bubbleserver"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// TODO: Move over to better spot
var serverStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

const banner = " /$$      /$$           /$$       /$$$$$$$            /$$       /$$       /$$                    \n| $$  /$ | $$          | $$      | $$__  $$          | $$      | $$      | $$                    \n| $$ /$$$| $$  /$$$$$$ | $$$$$$$ | $$  \\ $$ /$$   /$$| $$$$$$$ | $$$$$$$ | $$  /$$$$$$   /$$$$$$$\n| $$/$$ $$ $$ /$$__  $$| $$__  $$| $$$$$$$ | $$  | $$| $$__  $$| $$__  $$| $$ /$$__  $$ /$$_____/\n| $$$$_  $$$$| $$$$$$$$| $$  \\ $$| $$__  $$| $$  | $$| $$  \\ $$| $$  \\ $$| $$| $$$$$$$$|  $$$$$$ \n| $$$/ \\  $$$| $$_____/| $$  | $$| $$  \\ $$| $$  | $$| $$  | $$| $$  | $$| $$| $$_____/ \\____  $$\n| $$/   \\  $$|  $$$$$$$| $$$$$$$/| $$$$$$$/|  $$$$$$/| $$$$$$$/| $$$$$$$/| $$|  $$$$$$$ /$$$$$$$/\n|__/     \\__/ \\_______/|_______/ |_______/  \\______/ |_______/ |_______/ |__/ \\_______/|_______/ \n                                                                                                 \n                                                                                                 \n                                                                                                 \n"

const debug = true

type Tcommand struct {
	function func(string) string
	cmd      string
}

type (
	errMsg error
	screen string
)

var port = -1
var portString string

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

var commands = []Tcommand{}

const gap = "\n\n"

// Utility =============================================================================================================

func getDebugStatus(m Model) string {
	if debug {
		return m.statusStyle.Render("Status:") + " " + string(m.state) + "\n"
	} else {
		return ""
	}

}

// BubbleTea ===========================================================================================================

func main() {
	//Commands Setup
	addCommand(Tcommand{
		function: startServer,
		cmd:      "start server",
	})
	addCommand(Tcommand{
		function: stopServer,
		cmd:      "stop server",
	})
	addCommand(Tcommand{
		function: restartServer,
		cmd:      "restart server",
	})
	addCommand(Tcommand{
		function: setPort,
		cmd:      "set port",
	})

	// BubbleTea Setup ============
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	//Web setup
	err = web.WebInit()
	if err != nil {
		log.Fatal(err)
	}
	// Pass the program to the web server so it can send events
	web.SetProgram(p)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
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

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// Handle portSelection state separately - needs to process all input for the form
	if m.state == "portSelection" {
		// Update the form with the message
		form, cmd := m.portForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.portForm = f
		}

		// Check if form is complete
		if m.portForm.State == huh.StateCompleted {
			// Parse the port from portString
			if parsedPort, err := strconv.Atoi(portString); err == nil {
				port = parsedPort
				m.messages = append(m.messages,
					m.systemStyle.Render(fmt.Sprintf("Port set to %d", port)))
			} else {
				m.messages = append(m.messages,
					m.systemStyle.Render("Error: Invalid port number"))
			}

			// Return to chatting state
			m.state = "chatting"
			m.textarea.Focus()
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			// Reset form for next time
			m.portForm = generateForm()
		}

		return m, cmd
	}

	// Only update textarea and viewport if we're in chatting state
	if m.state == "chatting" {
		m.textarea, tiCmd = m.textarea.Update(msg)
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	switch msg := msg.(type) {
	case web.HTTPRequestMsg:
		// Handle HTTP request events from the web server
		m.messages = append(m.messages,
			serverStyle.Render("Server: ")+
				fmt.Sprintf("Received %s request to %s", msg.Method, msg.Path))
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		m.width = msg.Width
		m.height = msg.Height

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.state != "done" {
				m.messages = append(m.messages, m.systemStyle.Render("Please type \"done\" to quit"))
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
				m.textarea.Reset()
				m.viewport.GotoBottom()
			} else {
				fmt.Println(m.textarea.Value())
				return m, tea.Quit
			}

		case tea.KeyEnter:
			var textAreaVal = m.textarea.Value()

			if textAreaVal == "done" {
				m.state = "done"
				return m, m.spinner.Tick
			} else if textAreaVal == "/help" { //TODO Migrate this to a command
				m.messages = append(m.messages, m.systemStyle.Render("/start server - Start the web server\n/stop server - Stop the web server\n/restart server - Restart the web server\n/set port - Set the server port\n/help - Show this help message"))
				textAreaVal = ""
			} else if textAreaVal == "help" {
				m.messages = append(m.messages, m.systemStyle.Render("Did you mean /help?"))
				textAreaVal = ""
			} else {
				// For any messages
				//TODO Add in unknown handling
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			}

			//Command Handling
			if strings.HasPrefix(textAreaVal, "/") {
				textAreaVal = strings.TrimPrefix(textAreaVal, "/")
				result := HandleCommand(textAreaVal)

				// Check if it's a state change command
				if strings.HasPrefix(result, "STATE_CHANGE:") {
					newState := strings.TrimPrefix(result, "STATE_CHANGE:")
					m.state = screen(newState)
					m.textarea.Blur()
					return m, m.portForm.Init()
				}

				m.messages = append(m.messages, result)
			}

			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil

	default:
		if m.state == "done" {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	switch m.state {
	case "chatting":
		return fmt.Sprintf(
			"%s%s%s%s",
			getDebugStatus(m),
			m.viewport.View(),
			gap,
			m.textarea.View(),
		)
	case "portSelection":
		doc := fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			getDebugStatus(m),
			m.portStyle.Render("Set Server Port"),
			m.portForm.View(),
		)
		centered := lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			doc,
		)
		return centered
	case "done":
		return fmt.Sprintf(
			"%s%sWaiting for you to exit\n%s",
			getDebugStatus(m),
			m.spinner.View(),
			"Thank you! Press Ctrl+C or Esc to exit.",
		)
	}
	return "Unknown state"
}

// Command Handling ====================================================================================================

func HandleCommand(message string) string {
	for _, command := range commands {
		if message == command.cmd {
			return command.function(message)
		}
	}
	return fmt.Sprintf("\"/%s\" is not a valid command", message)
}

func startServer(string) string {
	return web.StartWebserverGoRoutine(port)
}
func stopServer(string) string    { return web.StopWebserver() }
func restartServer(string) string { return web.RestartWebserver(port) }
func setPort(string) string {
	return "STATE_CHANGE:portSelection"
}

func addCommand(newCommand Tcommand) {
	commands = append(commands, newCommand)
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
