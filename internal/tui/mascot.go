package tui

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	chafa "github.com/ploMP4/chafa-go"
	"github.com/charmbracelet/lipgloss"
)

// alturaMascote define quantas linhas de células o mascote ocupa.
// cols é calculado automaticamente preservando aspect ratio da imagem
// compensando a proporção 2:1 (altura:largura) das células do terminal.
const alturaMascote = 16

// renderMascote tenta usar chafa-go com assets/castor.png.
// Se o arquivo não existir, usa o fallback em lipgloss.
func renderMascote() string {
	if art, err := renderComChafa("assets/castor.png"); err == nil {
		return art
	}
	return renderMascoteFallback()
}

func renderComChafa(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// células do terminal têm ~2:1 (h:w), então cols = rows * 2 * (w/h)
	// para imagem quadrada: cols = rows * 2
	rows := alturaMascote
	cols := rows * 2 * w / h
	if cols < 1 {
		cols = rows * 2
	}

	// extrai bytes RGBA8
	pixels := make([]uint8, w*h*4)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			i := (y*w + x) * 4
			pixels[i], pixels[i+1], pixels[i+2], pixels[i+3] = uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		}
	}

	sm := chafa.SymbolMapNew()
	defer chafa.SymbolMapUnref(sm)
	chafa.SymbolMapAddByTags(sm, chafa.CHAFA_SYMBOL_TAG_VHALF)

	cfg := chafa.CanvasConfigNew()
	defer chafa.CanvasConfigUnref(cfg)
	chafa.CanvasConfigSetGeometry(cfg, int32(cols), int32(rows))
	chafa.CanvasConfigSetSymbolMap(cfg, sm)

	canvas := chafa.CanvasNew(cfg)
	defer chafa.CanvasUnRef(canvas)

	chafa.CanvasDrawAllPixels(
		canvas,
		chafa.CHAFA_PIXEL_RGBA8_UNASSOCIATED,
		pixels, int32(w), int32(h), int32(w*4),
	)

	return chafa.CanvasPrint(canvas, nil).String(), nil
}

// --- fallback lipgloss (usado quando assets/castor.png não existe) ---

type pxColor = lipgloss.Color

const (
	pxVazio  pxColor = ""
	pxEscuro pxColor = "#3D1F0A"
	pxMedio  pxColor = "#7B4F2E"
	pxClaro  pxColor = "#A57850"
	pxCreme  pxColor = "#F0DEB0"
	pxPreto  pxColor = "#111111"
	pxBranco pxColor = "#EEEEEE"
	pxRosa   pxColor = "#E89080"
	pxBege   pxColor = "#D4A574"
)

var castorGrid = [][]pxColor{
	{pxVazio, pxVazio, pxEscuro, pxMedio, pxMedio, pxMedio, pxMedio, pxEscuro, pxVazio, pxVazio},
	{pxVazio, pxEscuro, pxMedio, pxClaro, pxClaro, pxClaro, pxClaro, pxMedio, pxEscuro, pxVazio},
	{pxEscuro, pxMedio, pxClaro, pxClaro, pxClaro, pxClaro, pxClaro, pxClaro, pxMedio, pxEscuro},
	{pxEscuro, pxMedio, pxPreto, pxPreto, pxClaro, pxClaro, pxPreto, pxPreto, pxMedio, pxEscuro},
	{pxEscuro, pxMedio, pxPreto, pxBranco, pxClaro, pxClaro, pxPreto, pxBranco, pxMedio, pxEscuro},
	{pxEscuro, pxBege, pxCreme, pxCreme, pxCreme, pxCreme, pxCreme, pxCreme, pxBege, pxEscuro},
	{pxEscuro, pxCreme, pxCreme, pxRosa, pxRosa, pxRosa, pxRosa, pxCreme, pxCreme, pxEscuro},
	{pxEscuro, pxCreme, pxCreme, pxPreto, pxCreme, pxCreme, pxPreto, pxCreme, pxCreme, pxEscuro},
	{pxVazio, pxEscuro, pxMedio, pxCreme, pxCreme, pxCreme, pxCreme, pxMedio, pxEscuro, pxVazio},
	{pxVazio, pxEscuro, pxMedio, pxCreme, pxCreme, pxCreme, pxCreme, pxMedio, pxEscuro, pxVazio},
	{pxVazio, pxEscuro, pxMedio, pxMedio, pxCreme, pxCreme, pxMedio, pxMedio, pxEscuro, pxVazio},
	{pxVazio, pxVazio, pxEscuro, pxMedio, pxVazio, pxVazio, pxMedio, pxEscuro, pxVazio, pxVazio},
	{pxVazio, pxVazio, pxEscuro, pxEscuro, pxVazio, pxVazio, pxEscuro, pxEscuro, pxVazio, pxVazio},
}

var bigodeLinhas = map[int]bool{5: true, 6: true}

func renderMascoteFallback() string {
	var sb strings.Builder
	for row, linha := range castorGrid {
		if bigodeLinhas[row] {
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
			sb.WriteString(lipgloss.NewStyle().Foreground(pxEscuro).Render("━━"))
		}
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
