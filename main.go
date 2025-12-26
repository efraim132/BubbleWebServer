package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	web "BubbleWebServer/WebServer"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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
)

type Model struct {
	state       string
	spinner     spinner.Model
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	statusStyle lipgloss.Style
	systemStyle lipgloss.Style
	err         error
}

var commands = []Tcommand{}

const gap = "\n\n"

// Utility =============================================================================================================

func getDebugStatus(m Model) string {
	if debug {
		return m.statusStyle.Render("Status:") + " " + m.state + "\n"
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

	// BubbleTea Setup ============
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

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

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

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
				m.messages = append(m.messages, m.systemStyle.Render("/start<or>stop server\n/restart server\n/help to show this page"))
				textAreaVal = ""
			} else if textAreaVal == "help" {
				m.messages = append(m.messages, m.systemStyle.Render("Did you mean /help?"))
				textAreaVal = ""
			} else {
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			}

			//Command Handling
			if strings.HasPrefix(textAreaVal, "/") {
				textAreaVal = strings.TrimPrefix(textAreaVal, "/")
				m.messages = append(m.messages, HandleCommand(textAreaVal))
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
	if m.state == "chatting" {
		return fmt.Sprintf(
			"%s%s%s%s",
			getDebugStatus(m),
			m.viewport.View(),
			gap,
			m.textarea.View(),
		)
	} else if m.state == "done" {
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

func startServer(string) string   { return web.StartWebserverGoRoutine() }
func stopServer(string) string    { return web.StopWebserver() }
func restartServer(string) string { return web.RestartWebserver() }

func addCommand(newCommand Tcommand) error {
	commands = append(commands, newCommand)
	return nil
}
