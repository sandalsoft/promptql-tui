package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#06B6D4")
	mutedColor     = lipgloss.Color("#6B7280")
	errorColor     = lipgloss.Color("#EF4444")
	successColor   = lipgloss.Color("#10B981")
	bgColor        = lipgloss.Color("#1F2937")
	fgColor        = lipgloss.Color("#F9FAFB")
	borderColor    = lipgloss.Color("#374151")

	// App chrome
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Error display
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor)

	// Chat messages
	userMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			Bold(true)

	assistantMsgStyle = lipgloss.NewStyle().
				Foreground(fgColor)

	// List items
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(fgColor)

	// Input
	promptStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// Help text
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Borders
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2)
)
