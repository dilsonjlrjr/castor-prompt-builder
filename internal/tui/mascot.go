package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// cada "pixel" = 2 espaços com cor de fundo
type px = lipgloss.Color

const (
	pxVazio  px = ""        // transparente
	pxEscuro px = "#3D1F0A" // marrom escuro (contorno)
	pxMedio  px = "#7B4F2E" // marrom médio (corpo)
	pxClaro  px = "#A57850" // marrom claro (pelo)
	pxCreme  px = "#F0DEB0" // creme (focinho/barriga)
	pxPreto  px = "#111111" // olhos
	pxBranco px = "#EEEEEE" // brilho dos olhos
	pxRosa   px = "#E89080" // nariz/boca
	pxBege   px = "#D4A574" // bochecha
)

// grade 10×13 — inspirada no pixel-art de referência
var castorGrid = [][]px{
	// orelhas e topo da cabeça
	{pxVazio, pxVazio, pxEscuro, pxMedio, pxMedio, pxMedio, pxMedio, pxEscuro, pxVazio, pxVazio},
	{pxVazio, pxEscuro, pxMedio, pxClaro, pxClaro, pxClaro, pxClaro, pxMedio, pxEscuro, pxVazio},
	// testa
	{pxEscuro, pxMedio, pxClaro, pxClaro, pxClaro, pxClaro, pxClaro, pxClaro, pxMedio, pxEscuro},
	// olhos
	{pxEscuro, pxMedio, pxPreto, pxPreto, pxClaro, pxClaro, pxPreto, pxPreto, pxMedio, pxEscuro},
	{pxEscuro, pxMedio, pxPreto, pxBranco, pxClaro, pxClaro, pxPreto, pxBranco, pxMedio, pxEscuro},
	// bochecha + bigodes
	{pxEscuro, pxBege, pxCreme, pxCreme, pxCreme, pxCreme, pxCreme, pxCreme, pxBege, pxEscuro},
	// nariz e boca
	{pxEscuro, pxCreme, pxCreme, pxRosa, pxRosa, pxRosa, pxRosa, pxCreme, pxCreme, pxEscuro},
	{pxEscuro, pxCreme, pxCreme, pxPreto, pxCreme, pxCreme, pxPreto, pxCreme, pxCreme, pxEscuro},
	// corpo
	{pxVazio, pxEscuro, pxMedio, pxCreme, pxCreme, pxCreme, pxCreme, pxMedio, pxEscuro, pxVazio},
	{pxVazio, pxEscuro, pxMedio, pxCreme, pxCreme, pxCreme, pxCreme, pxMedio, pxEscuro, pxVazio},
	{pxVazio, pxEscuro, pxMedio, pxMedio, pxCreme, pxCreme, pxMedio, pxMedio, pxEscuro, pxVazio},
	// pés
	{pxVazio, pxVazio, pxEscuro, pxMedio, pxVazio, pxVazio, pxMedio, pxEscuro, pxVazio, pxVazio},
	{pxVazio, pxVazio, pxEscuro, pxEscuro, pxVazio, pxVazio, pxEscuro, pxEscuro, pxVazio, pxVazio},
}

// bigodes — linhas que recebem extensões laterais
var bigodeLinhas = map[int]bool{5: true, 6: true}

func renderMascote() string {
	var sb strings.Builder
	for row, linha := range castorGrid {
		if bigodeLinhas[row] {
			// bigode esquerdo
			sb.WriteString(lipgloss.NewStyle().Foreground(pxEscuro).Render("━━"))
		} else {
			sb.WriteString("  ")
		}
		for _, cor := range linha {
			if cor == pxVazio {
				sb.WriteString("  ")
			} else {
				sb.WriteString(lipgloss.NewStyle().Background(cor).Render("  "))
			}
		}
		if bigodeLinhas[row] {
			// bigode direito
			sb.WriteString(lipgloss.NewStyle().Foreground(pxEscuro).Render("━━"))
		}
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
