package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// TUI Integration =====================================================================================================
// This file provides a Terminal User Interface (TUI) for the server
// When you use StartWithTUI(), instead of just logging to the console,
// you get a nice graphical interface in your terminal showing server status and requests

// TUI Configuration ===================================================================================================

// TUIConfig holds configuration for the TUI
type TUIConfig struct {
	Banner      string // ASCII art banner shown at startup
	Title       string // Title shown in the TUI
	ShowDebug   bool   // Whether to show debug information
	BannerColor string // Color for the banner (lipgloss color)
	ServerColor string // Color for server messages (lipgloss color)
	SystemColor string // Color for system messages (lipgloss color)
}

// DefaultTUIConfig returns a default TUI configuration
func DefaultTUIConfig() TUIConfig {
	return TUIConfig{
		Banner: `
 __        __   _      ____
 \ \      / /__| |__  / ___|  ___ _ ____   _____ _ __
  \ \ /\ / / _ \ '_ \ \___ \ / _ \ '__\ \ / / _ \ '__|
   \ V  V /  __/ |_) | ___) |  __/ |   \ V /  __/ |
    \_/\_/ \___|_.__/ |____/ \___|_|    \_/ \___|_|
`,
		Title:       "Web Server Control Panel",
		ShowDebug:   false,
		BannerColor: "248",
		ServerColor: "5",
		SystemColor: "1",
	}
}

// TUI Model and State =================================================================================================

// tuiModel represents the state of the TUI
// This is what BubbleTea uses to track the current state of the interface
type tuiModel struct {
	server      *Server        // Reference to the server we're controlling
	viewport    viewport.Model // The scrollable area showing messages
	textarea    textarea.Model // The input field for commands
	messages    []string       // All messages displayed in the viewport
	config      TUIConfig      // TUI configuration
	bannerStyle lipgloss.Style // Style for the banner
	serverStyle lipgloss.Style // Style for server messages
	systemStyle lipgloss.Style // Style for system messages
	statusStyle lipgloss.Style // Style for status line
	portStyle   lipgloss.Style // Style for port selection form
	width       int            // Current terminal width
	height      int            // Current terminal height
	quitting    bool           // Whether we're in the process of quitting
	state       string         // Current TUI state ("chatting" or "portSelection")
	portForm    *huh.Form      // Form for port selection
	portString  string         // String value from port form
}

// HTTPRequestMsg is a message sent when an HTTP request is received
// BubbleTea works by sending messages - this is how the server tells the TUI about requests
type HTTPRequestMsg struct {
	Method string // HTTP method (GET, POST, etc.)
	Path   string // URL path that was requested
	IP     string // Client IP address
}

// TUI Initialization ==================================================================================================

// newTUIModel creates a new TUI model with the given configuration
func newTUIModel(srv *Server, config TUIConfig) tuiModel {
	// Create the text input area for commands
	ta := textarea.New()
	ta.Placeholder = "Enter command..."
	ta.Focus()
	ta.Prompt = "> "
	ta.CharLimit = 280
	ta.SetWidth(30)
	ta.SetHeight(1)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false) // Enter sends command instead of newline

	// Create the viewport (scrollable message area)
	vp := viewport.New(30, 5)

	// Build the welcome message
	welcomeMsg := config.Banner + "\n" + config.Title + "\n\n"
	welcomeMsg += "Server ready on port " + fmt.Sprintf("%d", srv.port) + "\n"
	welcomeMsg += "Commands: /start, /stop, /restart, /port, /help, quit\n"

	vp.SetContent(welcomeMsg)

	// Create styles
	bannerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(config.BannerColor))
	serverStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(config.ServerColor))
	systemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(config.SystemColor))
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	portStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	// Create the model
	model := tuiModel{
		server:      srv,
		viewport:    vp,
		textarea:    ta,
		messages:    []string{},
		config:      config,
		bannerStyle: bannerStyle,
		serverStyle: serverStyle,
		systemStyle: systemStyle,
		statusStyle: statusStyle,
		portStyle:   portStyle,
		width:       30,
		height:      5,
		quitting:    false,
		state:       "chatting",
		portString:  "",
	}

	// Generate the port form
	model.portForm = model.generatePortForm()

	return model
}

// BubbleTea Interface Methods =========================================================================================
// These methods implement the BubbleTea Model interface

// Init initializes the TUI
func (m tuiModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages and updates the model
// This is called whenever something happens (key press, message received, etc.)
func (m *tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// Handle port selection state separately - form needs to process all input
	if m.state == "portSelection" {
		// Update the form with the message
		form, cmd := m.portForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.portForm = f
		}

		// Check if form is complete
		if m.portForm.State == huh.StateCompleted {
			// Trim whitespace from the port string
			portStr := strings.TrimSpace(m.portString)

			// Parse the port from portString
			if parsedPort, err := strconv.Atoi(portStr); err == nil {
				// Check if server is running - can't change port while running
				if m.server.IsRunning() {
					m.messages = append(m.messages,
						m.systemStyle.Render("System: Cannot change port while server is running. Stop the server first."))
				} else {
					// Update the server's port
					m.server.port = parsedPort
					m.messages = append(m.messages,
						m.systemStyle.Render(fmt.Sprintf("System: Port set to %d", parsedPort)))
				}
			} else {
				// Debug: show what we received
				m.messages = append(m.messages,
					m.systemStyle.Render(fmt.Sprintf("System: Error - Invalid port number (received: '%s', error: %v)", m.portString, err)))
			}

			// Return to chatting state
			m.state = "chatting"
			m.textarea.Focus()
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			// Reset form for next time
			m.portForm = m.generatePortForm()
		}

		return m, cmd
	}

	// Only update textarea and viewport if we're in chatting state
	if m.state == "chatting" {
		m.textarea, tiCmd = m.textarea.Update(msg)
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	switch msg := msg.(type) {
	case HTTPRequestMsg:
		// An HTTP request was received - log it in the viewport
		m.messages = append(m.messages,
			m.serverStyle.Render("Server: ")+
				fmt.Sprintf("Received %s request to %s from %s", msg.Method, msg.Path, msg.IP))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

	case tea.WindowSizeMsg:
		// Terminal was resized - adjust our layout
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - 4
		m.width = msg.Width
		m.height = msg.Height

		if len(m.messages) > 0 {
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
		}
		m.viewport.GotoBottom()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// Ctrl+C or Esc - quit
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			// Enter - process the command
			command := m.textarea.Value()
			m.textarea.Reset()

			if command == "" {
				return m, tea.Batch(tiCmd, vpCmd)
			}

			// Handle commands
			response := m.handleCommand(command)
			if response != "" {
				m.messages = append(m.messages, response)
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			}

			// Check if state changed to portSelection - need to init the form
			if m.state == "portSelection" {
				return m, m.portForm.Init()
			}

			// Check if we should quit
			if command == "quit" || command == "exit" {
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

// View renders the TUI
// This is called whenever the screen needs to be redrawn
func (m *tuiModel) View() string {
	if m.quitting {
		return "Shutting down server...\n"
	}

	// Debug status line (optional)
	debug := ""
	if m.config.ShowDebug {
		debug = m.statusStyle.Render(fmt.Sprintf("Debug: %dx%d | Messages: %d | State: %s\n", m.width, m.height, len(m.messages), m.state))
	}

	// Handle different states
	switch m.state {
	case "chatting":
		// Normal chat interface
		return fmt.Sprintf(
			"%s%s\n\n%s",
			debug,
			m.viewport.View(),
			m.textarea.View(),
		)

	case "portSelection":
		// Port selection form (centered)
		doc := fmt.Sprintf(
			"%s\n%s\n\n%s",
			debug,
			m.portStyle.Render("Set Server Port"),
			m.portForm.View(),
		)
		// Center the form on the screen
		centered := lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			doc,
		)
		return centered

	default:
		return "Unknown state"
	}
}

// Command Handling ====================================================================================================

// handleCommand processes a command entered by the user
func (m *tuiModel) handleCommand(cmd string) string {
	// Remove leading slash if present
	cmd = strings.TrimPrefix(cmd, "/")

	switch cmd {
	case "start":
		if m.server.IsRunning() {
			return m.systemStyle.Render("System: Server is already running")
		}
		if err := m.server.Start(); err != nil {
			return m.systemStyle.Render(fmt.Sprintf("System: Failed to start server: %v", err))
		}
		return m.serverStyle.Render(fmt.Sprintf("Server: Started on port %d", m.server.GetPort()))

	case "stop":
		if !m.server.IsRunning() {
			return m.systemStyle.Render("System: Server is not running")
		}
		if err := m.server.Stop(); err != nil {
			return m.systemStyle.Render(fmt.Sprintf("System: Failed to stop server: %v", err))
		}
		return m.serverStyle.Render("Server: Stopped")

	case "restart":
		if !m.server.IsRunning() {
			return m.systemStyle.Render("System: Server is not running")
		}
		if err := m.server.Restart(); err != nil {
			return m.systemStyle.Render(fmt.Sprintf("System: Failed to restart server: %v", err))
		}
		return m.serverStyle.Render(fmt.Sprintf("Server: Restarted on port %d", m.server.GetPort()))

	case "status":
		if m.server.IsRunning() {
			return m.serverStyle.Render(fmt.Sprintf("Server: Running on %s:%d", m.server.GetHost(), m.server.GetPort()))
		}
		return m.serverStyle.Render("Server: Not running")

	case "port", "setport":
		// Transition to port selection state
		m.state = "portSelection"
		m.textarea.Blur()
		// Regenerate the form with the current model pointer
		m.portForm = m.generatePortForm()
		return ""

	case "help":
		return m.systemStyle.Render(`Commands:
  /start   - Start the server
  /stop    - Stop the server
  /restart - Restart the server
  /status  - Show server status
  /port    - Change server port
  /help    - Show this help
  quit     - Exit the TUI`)

	case "quit", "exit":
		// Handled in Update() to trigger tea.Quit
		return ""

	default:
		return m.systemStyle.Render(fmt.Sprintf("Unknown command: /%s (type /help for commands)", cmd))
	}
}

// Port Selection Form =================================================================================================

// generatePortForm creates a new Huh form for port selection
// This creates an interactive form that validates port numbers
func (m *tuiModel) generatePortForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("HTTP Port").
				Placeholder("8080").
				Value(&m.portString). // Pointer to string variable in the model
				// Validate function checks if it's a valid port number
				Validate(func(s string) error {
					// Allow empty string (user is still typing)
					if s == "" {
						return nil
					}

					// Check if it can be converted to integer
					port, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("must be a number")
					}

					// Check if port is in valid range (1-65535)
					if port < 1 || port > 65535 {
						return fmt.Errorf("must be between 1 and 65535")
					}

					// Valid port!
					return nil
				}),
		),
	).WithWidth(60) // Set the form width

	return form
}

// TUI Integration with Server =========================================================================================

// StartWithTUI starts the server with a Terminal User Interface
// This gives you a graphical interface in the terminal for controlling the server
// instead of just log statements
//
// Example:
//
//	srv := server.New().Port(8080).Route("/hello", handler).Build()
//	srv.StartWithTUI()  // Blocks until user quits the TUI
func (s *Server) StartWithTUI() error {
	return s.StartWithTUIConfig(DefaultTUIConfig())
}

// StartWithTUIConfig starts the server with a custom TUI configuration
// This allows you to customize the banner, colors, etc.
//
// Example:
//
//	config := server.DefaultTUIConfig()
//	config.Banner = "My Custom Banner"
//	config.ShowDebug = true
//	srv.StartWithTUIConfig(config)
func (s *Server) StartWithTUIConfig(config TUIConfig) error {
	// Set up the OnRequest callback to send messages to the TUI
	// We'll store the BubbleTea program in a variable so the callback can use it
	var program *tea.Program

	s.OnRequest(func(info RequestInfo) {
		if program != nil {
			program.Send(HTTPRequestMsg{
				Method: info.Method,
				Path:   info.Path,
				IP:     info.IP,
			})
		}
	})

	// Create the TUI model (as pointer since Update/View use pointer receivers)
	model := newTUIModel(s, config)

	// Create the BubbleTea program
	program = tea.NewProgram(&model, tea.WithAltScreen())

	// Run the TUI (this blocks until the user quits)
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("TUI error: %v", err)
	}

	// When the TUI exits, stop the server if it's running
	if s.IsRunning() {
		s.Stop()
	}

	return nil
}

// How This Works ======================================================================================================
//
// 1. You create a server and configure it normally:
//    srv := server.New().Port(8080).Route("/hello", handler).Build()
//
// 2. Instead of srv.Start(), you call srv.StartWithTUI():
//    srv.StartWithTUI()
//
// 3. This launches a TUI that shows:
//    - A welcome banner
//    - Server status
//    - Real-time HTTP request logs
//    - A command prompt for controlling the server
//
// 4. You can type commands like:
//    /start   - Start the server
//    /stop    - Stop the server
//    /restart - Restart the server
//    /port    - Change the server port (interactive form)
//    /status  - Show server status
//    /help    - Show all commands
//    quit     - Exit
//
// 5. Port Selection:
//    - Type /port to open an interactive form
//    - The form validates port numbers (must be 1-65535)
//    - Allows any valid port (OS will handle privilege checks)
//    - Cannot change port while server is running
//
// 6. When you quit, the server automatically stops
//
// The TUI uses BubbleTea, which is a framework for building terminal UIs
// It works by:
//   - Model: Holds the state (messages, server reference, etc.)
//   - Update: Handles events (key presses, messages from server)
//   - View: Renders the UI
//
// Messages flow like this:
//   HTTP Request → Server → OnRequest callback → HTTPRequestMsg → TUI Update → TUI View
