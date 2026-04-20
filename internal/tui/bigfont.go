package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Cada letra tem 8 linhas de altura × 6 colunas de largura.
// Duas palavras (CASTOR + BUILDER) = 8 + 1 espaço + 8 = 17 linhas ≈ alturaMascote.
const fontHeight = 8

var bigFont = map[rune][]string{
	'C': {
		"  ████",
		" ██   ",
		"██    ",
		"██    ",
		"██    ",
		"██    ",
		" ██   ",
		"  ████",
	},
	'A': {
		"  ██  ",
		" █  █ ",
		"██  ██",
		"██  ██",
		"██████",
		"██  ██",
		"██  ██",
		"██  ██",
	},
	'S': {
		" █████",
		"██    ",
		"██    ",
		" ████ ",
		"    ██",
		"    ██",
		"    ██",
		"█████ ",
	},
	'T': {
		"██████",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
	},
	'O': {
		" ████ ",
		"██  ██",
		"██  ██",
		"██  ██",
		"██  ██",
		"██  ██",
		"██  ██",
		" ████ ",
	},
	'R': {
		"█████ ",
		"██  ██",
		"██  ██",
		"█████ ",
		"██ ██ ",
		"██  ██ ",
		"██  ██",
		"██  ██",
	},
	'B': {
		"█████ ",
		"██  ██",
		"██  ██",
		"█████ ",
		"██  ██",
		"██  ██",
		"██  ██",
		"█████ ",
	},
	'U': {
		"██  ██",
		"██  ██",
		"██  ██",
		"██  ██",
		"██  ██",
		"██  ██",
		"██  ██",
		" ████ ",
	},
	'I': {
		"██████",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"██████",
	},
	'L': {
		"██    ",
		"██    ",
		"██    ",
		"██    ",
		"██    ",
		"██    ",
		"██    ",
		"██████",
	},
	'D': {
		"████  ",
		"██  █ ",
		"██   █",
		"██   █",
		"██   █",
		"██   █",
		"██  █ ",
		"████  ",
	},
	'E': {
		"██████",
		"██    ",
		"██    ",
		"████  ",
		"██    ",
		"██    ",
		"██    ",
		"██████",
	},
	' ': {
		"   ",
		"   ",
		"   ",
		"   ",
		"   ",
		"   ",
		"   ",
		"   ",
	},
}

// renderPalavraGrande renderiza uma palavra usando o bigFont, com a cor especificada.
func renderPalavraGrande(palavra string, cor lipgloss.Color) string {
	rows := make([]string, fontHeight)
	st := lipgloss.NewStyle().Foreground(cor).Bold(true)

	for _, ch := range strings.ToUpper(palavra) {
		glyph, ok := bigFont[ch]
		if !ok {
			glyph = bigFont[' ']
		}
		for i := range rows {
			rows[i] += glyph[i] + " "
		}
	}

	var sb strings.Builder
	for _, row := range rows {
		sb.WriteString(st.Render(row) + "\n")
	}
	return sb.String()
}

// bigTitle retorna "CASTOR" e "BUILDER" em blocos grandes, alinhados à altura do mascote.
func bigTitle() string {
	sep := lipgloss.NewStyle().Foreground(colorMuted).Render(strings.Repeat("─", 44))
	return renderPalavraGrande("CASTOR", colorPrimary) +
		sep + "\n" +
		renderPalavraGrande("BUILDER", colorAccent)
}
