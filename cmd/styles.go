package cmd

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles contains all the lipgloss styles for the tt CLI.
// This file defines a consistent color scheme and styling approach to make
// the CLI look professional and modern. The theme uses a blue/green primary
// palette suitable for a git helper tool, with clear distinctions for
// success, error, info, and highlight states.
//
// Color Scheme:
// - Primary: Blue (#00BFFF) for headers and main elements
// - Success: Green (#00FF7F) for positive feedback (e.g., branch created, commit successful)
// - Error: Red (#FF4500) for failures and warnings
// - Info: Teal (#20B2AA) for informational messages
// - Highlight: Yellow (#FFD700) for important elements like branch names
// - Neutral: Gray (#808080) for subtle text and borders
//
// Usage:
// - Import this package in command files: import "github.com/aixoio/tt/cmd"
// - Access styles via Styles.Success.Render("message") for styled output
// - For headers: Styles.Header.Render("Title")
// - For prompts: Use with huh or directly for custom outputs
// - Always use lipgloss.Join or NewRenderer for complex layouts
// - Ensure terminal supports 256 colors or truecolor for best results
//
// Future Integration:
// - Replace fmt.Printf with styled renders in all commands (branch.go, commit.go, etc.)
// - Style huh inputs with custom themes matching this scheme
// - Add icons: ✓ for success, ✗ for errors, ℹ for info
// - Consistent padding/margins: 1 space around messages

var (
	// Primary styles
	Primary = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF")).
		Bold(true)

	Header = Primary.Underline(true).MarginBottom(1)

	// Success styles
	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF7F")).
		Bold(true)

	SuccessIcon = Success.Render("✓ ")

	// Error styles
	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF4500")).
		Bold(true)

	ErrorIcon = Error.Render("✗ ")

	// Info styles
	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#20B2AA"))

	InfoIcon = Info.Render("ℹ ")

	// Highlight styles
	Highlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	// Neutral styles
	Neutral = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080"))

	// Block styles for sections
	Block = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#808080")).
		Padding(1)

	// Inline styles for commands
	InlineCode = lipgloss.NewStyle().
			Background(lipgloss.Color("#1E1E1E")).
			Foreground(lipgloss.Color("#D4D4D4")).
			Padding(0, 1)
)
