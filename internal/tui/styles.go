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

	// cores do castor pixel-art
	colorCastor     = lipgloss.Color("#7B3F00") // marrom escuro
	colorCastorLight = lipgloss.Color("#A0522D") // marrom claro
	colorTeeth      = lipgloss.Color("#FFFFF0")  // dentes marfim

	styleCastor = lipgloss.NewStyle().
			Foreground(colorCastorLight).
			Background(lipgloss.NoColor{})

	styleCastorAccent = lipgloss.NewStyle().
				Foreground(colorTeeth).
				Bold(true)

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

// castorFrames — 4 quadros pixel-art do rosto do castor dançando
// inspirado em pixel art: orelhas, olhos brancos, nariz, dentes
var castorFrames = []string{
	// frame 0 — neutro
	" ▄█▄ ▄█▄ \n" +
		"███████████\n" +
		"█▐██▌▐██▌█\n" +
		"█▐██▌▐██▌█\n" +
		"█   ▄▄   █\n" +
		"█  █  █  █\n" +
		" █████████\n" +
		"  ▐█▌▐█▌  ",

	// frame 1 — piscando olho esquerdo (dança)
	" ▄█▄ ▄█▄ \n" +
		"███████████\n" +
		"█▐──▌▐██▌█\n" +
		"█▐──▌▐██▌█\n" +
		"█   ▄▄   █\n" +
		"█  █  █  █\n" +
		" █████████\n" +
		"  ▐█▌▐█▌  ",

	// frame 2 — olhos arregalados (dança)
	" ▄█▄ ▄█▄ \n" +
		"███████████\n" +
		"█▐██▌▐██▌█\n" +
		"█▐██▌▐██▌█\n" +
		"█   ▄▄   █\n" +
		"█ █    █ █\n" +
		" █████████\n" +
		"  ▐█▌▐█▌  ",

	// frame 3 — feliz/pulando (★ nos olhos)
	" ▄█▄ ▄█▄ \n" +
		"███████████\n" +
		"█ ★    ★  █\n" +
		"█         █\n" +
		"█   ▄▄   █\n" +
		"█ ██  ██ █\n" +
		" █████████\n" +
		"  ▐█▌▐█▌  ",
}
