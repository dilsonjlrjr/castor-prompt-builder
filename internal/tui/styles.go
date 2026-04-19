package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorPrimary   = lipgloss.Color("#F4A261") // laranja castor
	colorSecondary = lipgloss.Color("#2A9D8F")
	colorMuted     = lipgloss.Color("#6B7280")
	colorAccent    = lipgloss.Color("#E9C46A")
	colorBg        = lipgloss.Color("#1A1A2E")
	colorText      = lipgloss.Color("#E2E8F0")
	colorError     = lipgloss.Color("#EF4444")
	colorSuccess   = lipgloss.Color("#10B981")

	styleTitle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			PaddingLeft(1)

	styleSubtitle = lipgloss.NewStyle().
			Foreground(colorSecondary).
			PaddingLeft(1)

	styleHelp = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(1)

	styleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	styleSelected = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	styleNormal = lipgloss.NewStyle().
			Foreground(colorText)

	styleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleSuccess = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	styleError = lipgloss.NewStyle().
			Foreground(colorError)

	styleHeader = lipgloss.NewStyle().
			Background(colorPrimary).
			Foreground(colorBg).
			Bold(true).
			Padding(0, 2)

	styleBadge = lipgloss.NewStyle().
			Background(colorSecondary).
			Foreground(colorBg).
			Padding(0, 1)
)

