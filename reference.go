package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// temp
var serverStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

const debug = false

func getDebugStatus(m Model) string {
	if debug {
		return m.statusStyle.Render("Status:") + " " + m.state + "\n"
	} else {
		return ""
	}

}

const gap = "\n\n"

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
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
	err         error
}

func initialModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		state:       "chatting",
		spinner:     spinner.New(spinner.WithSpinner(spinner.Dot)),
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		statusStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
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
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			var textAreaVal = m.textarea.Value()
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			if textAreaVal == "done" {
				m.state = "done"
				return m, m.spinner.Tick
			} else if textAreaVal == "server" {
				m.messages = append(m.messages, serverStyle.Render("Server: ")+"Hello World!")
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
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
			"Thank you for chatting! Press Ctrl+C or Esc to exit.",
		)
	}
	return "Unknown state"
}
