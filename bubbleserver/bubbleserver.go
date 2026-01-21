package bubbleserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

// Config holds the configuration for the BubbleServer
type Config struct {
	Banner      string
	WelcomeText string
	Debug       bool
	LogFile     string
}

// Command represents a TUI command that can be executed
type Command struct {
	Name        string
	Description string
	Handler     func(string) string
}

// HTTPRoute represents an HTTP endpoint
type HTTPRoute struct {
	Path    string
	Handler func(w interface{}, req interface{})
}

// Server manages the BubbleTea TUI and web server
type Server struct {
	config      Config
	commands    []Command
	httpRoutes  []HTTPRoute
	program     *tea.Program
	webServer   *webServerManager
	defaultPort int
}

type (
	errMsg error
	screen string
)

type model struct {
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
	port        int
	portString  string
	server      *Server
}

// NewServer creates a new BubbleServer instance
func NewServer(config Config) *Server {
	if config.LogFile == "" {
		config.LogFile = "debug.log"
	}
	if config.Banner == "" {
		config.Banner = " __      __      ___.  __________     ___.  ___.   .__     \n" +
			"/  \\    /  \\ ____\\_ |__\\______   \\__ _\\_ |__\\_ |__ |  |   ____   ______\n" +
			"\\   \\//\\   // __ \\| __ \\|    |  _/  |  \\ __ \\| __ \\|  | _/ __ \\ /  ___/\n" +
			" \\        /\\  ___/| \\_\\ \\    |   \\  |  / \\_\\ \\ \\_\\ \\  |_\\  ___/ \\___ \\ \n" +
			"  \\__/\\  /  \\___  >___  /______  /____/|___  /___  /____/\\___  >____  >\n" +
			"       \\/       \\/    \\/       \\/          \\/    \\/          \\/     \\/  \n"
	}
	if config.WelcomeText == "" {
		config.WelcomeText = "Welcome to the control panel!\nType a message and press Enter to send."
	}

	return &Server{
		config:      config,
		commands:    []Command{},
		httpRoutes:  []HTTPRoute{},
		defaultPort: 8080,
		webServer:   newWebServerManager(),
	}
}

// RegisterCommand adds a new TUI command
func (s *Server) RegisterCommand(cmd Command) {
	s.commands = append(s.commands, cmd)
}

// RegisterHTTPRoute adds a new HTTP route
func (s *Server) RegisterHTTPRoute(route HTTPRoute) {
	s.httpRoutes = append(s.httpRoutes, route)
}

// SetDefaultPort sets the default port for the web server
func (s *Server) SetDefaultPort(port int) {
	s.defaultPort = port
}

// Run starts the BubbleServer
func (s *Server) Run() error {
	// Setup default commands
	s.setupDefaultCommands()

	// Setup logging
	f, err := tea.LogToFile(s.config.LogFile, "debug")
	if err != nil {
		return fmt.Errorf("failed to setup logging: %w", err)
	}
	defer func() { _ = f.Close() }()

	// Initialize web server with routes
	if err := s.webServer.init(s.httpRoutes); err != nil {
		return fmt.Errorf("failed to initialize web server: %w", err)
	}

	// Create BubbleTea program
	s.program = tea.NewProgram(s.initialModel(), tea.WithAltScreen())

	// Set program reference in web server
	s.webServer.setProgram(s.program)

	// Run the program
	if _, err := s.program.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}

func (s *Server) setupDefaultCommands() {
	// Start server command
	s.RegisterCommand(Command{
		Name:        "start server",
		Description: "Start the web server",
		Handler: func(string) string {
			return s.webServer.start(s.defaultPort)
		},
	})

	// Stop server command
	s.RegisterCommand(Command{
		Name:        "stop server",
		Description: "Stop the web server",
		Handler: func(string) string {
			return s.webServer.stop()
		},
	})

	// Restart server command
	s.RegisterCommand(Command{
		Name:        "restart server",
		Description: "Restart the web server",
		Handler: func(string) string {
			return s.webServer.restart(s.defaultPort)
		},
	})

	// Set port command
	s.RegisterCommand(Command{
		Name:        "set port",
		Description: "Set the server port",
		Handler: func(string) string {
			return "STATE_CHANGE:portSelection"
		},
	})

	// Help command
	s.RegisterCommand(Command{
		Name:        "help",
		Description: "Show available commands",
		Handler: func(string) string {
			var helpText strings.Builder
			for _, cmd := range s.commands {
				helpText.WriteString(fmt.Sprintf("/%s - %s\n", cmd.Name, cmd.Description))
			}
			return helpText.String()
		},
	})
}

func (s *Server) initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Enter Command . . ."
	ta.Focus()
	ta.Prompt = "> "
	ta.CharLimit = 280
	ta.SetWidth(30)
	ta.SetHeight(1)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(s.config.Banner + "\n\n" + s.config.WelcomeText)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	m := model{
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
		bannerStyle: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}),
		width:       30,
		height:      5,
		err:         nil,
		port:        s.defaultPort,
		portString:  "",
		server:      s,
	}

	m.portForm = m.generateForm()
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// Handle portSelection state
	if m.state == "portSelection" {
		form, cmd := m.portForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.portForm = f
		}

		if m.portForm.State == huh.StateCompleted {
			if parsedPort, err := strconv.Atoi(m.portString); err == nil {
				m.port = parsedPort
				m.server.defaultPort = parsedPort
				m.messages = append(m.messages,
					m.systemStyle.Render(fmt.Sprintf("Port set to %d", m.port)))
			} else {
				m.messages = append(m.messages,
					m.systemStyle.Render("Error: Invalid port number"))
			}

			m.state = "chatting"
			m.textarea.Focus()
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			m.portForm = m.generateForm()
		}

		return m, cmd
	}

	// Update textarea and viewport in chatting state
	if m.state == "chatting" {
		m.textarea, tiCmd = m.textarea.Update(msg)
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	switch msg := msg.(type) {
	case HTTPRequestMsg:
		m.messages = append(m.messages,
			m.senderStyle.Render("Server: ")+
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
				return m, tea.Quit
			}

		case tea.KeyEnter:
			textAreaVal := m.textarea.Value()

			if textAreaVal == "done" {
				m.state = "done"
				return m, m.spinner.Tick
			}

			// Add user message
			if textAreaVal != "" {
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+textAreaVal)
			}

			// Handle commands
			if strings.HasPrefix(textAreaVal, "/") {
				textAreaVal = strings.TrimPrefix(textAreaVal, "/")
				result := m.handleCommand(textAreaVal)

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

func (m model) View() string {
	switch m.state {
	case "chatting":
		debugStatus := ""
		if m.server.config.Debug {
			debugStatus = m.statusStyle.Render("Status: ") + string(m.state) + "\n"
		}
		return fmt.Sprintf(
			"%s%s%s%s",
			debugStatus,
			m.viewport.View(),
			gap,
			m.textarea.View(),
		)

	case "portSelection":
		debugStatus := ""
		if m.server.config.Debug {
			debugStatus = m.statusStyle.Render("Status: ") + string(m.state) + "\n"
		}
		doc := fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			debugStatus,
			m.portStyle.Render("Set Server Port"),
			m.portForm.View(),
		)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, doc)

	case "done":
		return fmt.Sprintf(
			"%sWaiting for you to exit\n%s",
			m.spinner.View(),
			"Thank you! Press Ctrl+C or Esc to exit.",
		)
	}

	return "Unknown state"
}

func (m model) handleCommand(message string) string {
	for _, command := range m.server.commands {
		if message == command.Name {
			return command.Handler(message)
		}
	}
	return fmt.Sprintf("\"/%s\" is not a valid command. Type /help for available commands.", message)
}

func (m model) generateForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("HTTP Port").
				Placeholder("8080").
				Value(&m.portString).
				Validate(func(s string) error {
					port, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("port must be a number")
					}
					if port < 1 || port > 65535 {
						return fmt.Errorf("port must be between 1 and 65535")
					}
					if port < 1024 {
						return fmt.Errorf("port %d requires admin privileges", port)
					}
					return nil
				}),
		),
	).WithWidth(60)

	return form
}
