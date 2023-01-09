package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify"
)

/* CONSTANTS */

var (
	// P the current tea program
	P *tea.Program
	// client
	Client *spotify.Client
	// WindowSize store the size of the terminal window
	WindowSize tea.WindowSizeMsg
)

/* STYLING */

// DocStyle styling for viewports
var DocStyle = lipgloss.NewStyle().Margin(0, 2)

// HelpStyle styling for help context menu
var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

// ErrStyle provides styling for error messages
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd534b")).Render

// AlertStyle provides styling for alert messages
var AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render

type keymap struct {
	Radio  key.Binding
	Enter  key.Binding
	Rename key.Binding
	Delete key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Radio: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "create"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
