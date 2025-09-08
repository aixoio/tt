package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles contains all the lipgloss styles for the tt CLI.
// This file defines a consistent color scheme and styling approach to make
// the CLI look professional and modern. The theme uses a sophisticated
// color palette with excellent contrast and accessibility.
//
// Color Scheme:
// - Primary: Deep Blue (#2563EB) for headers and main elements
// - Success: Emerald Green (#10B981) for positive feedback
// - Error: Red (#EF4444) for failures and warnings
// - Info: Sky Blue (#0EA5E9) for informational messages
// - Warning: Amber (#F59E0B) for caution messages
// - Highlight: Violet (#8B5CF6) for important elements like branch names
// - Neutral: Cool Gray (#6B7280) for subtle text
// - Muted: Warm Gray (#9CA3AF) for less important text
// - Border: Light Gray (#E5E7EB) for borders and dividers

var (
	// Primary styles - Deep blue for main elements
	Primary = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2563EB")).
		Bold(true)

	Header = Primary.
		Underline(true).
		MarginBottom(1)

	// Success styles - Emerald green for positive feedback
	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true)

	SuccessIcon = Success.Render("✓")

	// Error styles - Red for failures and warnings
	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Bold(true)

	ErrorIcon = Error.Render("✗")

	// Warning styles - Amber for caution messages
	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Bold(true)

	WarningIcon = Warning.Render("⚠")

	// Info styles - Sky blue for informational messages
	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0EA5E9"))

	InfoIcon = Info.Render("ℹ")

	// Highlight styles - Violet for important elements
	Highlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B5CF6")).
			Bold(true)

	// Neutral styles - Cool gray for subtle text
	Neutral = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	// Muted styles - Warm gray for less important text
	Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF"))

	// Block styles for sections
	Block = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E5E7EB")).
		Padding(1, 2).
		MarginBottom(1)

	// Card styles for important information
	Card = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2563EB")).
		Padding(1, 2).
		MarginBottom(1)

	// Inline styles for commands
	InlineCode = lipgloss.NewStyle().
			Background(lipgloss.Color("#F3F4F6")).
			Foreground(lipgloss.Color("#1F2937")).
			Padding(0, 1).
			Bold(true)

	// Separator style
	Separator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Render("─")

	// Spinner style for loading states
	Spinner = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B5CF6"))

	// Git command style
	GitCommand = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#059669")).
			Bold(true)

	// Branch style
	Branch = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7C3AED")).
		Bold(true)

	// Commit hash style
	CommitHash = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DC2626")).
			Bold(true)

	// File path style
	FilePath = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0891B2")).
			Italic(true)

	// Spinner icon
	SpinnerIcon = "⏳"

)
