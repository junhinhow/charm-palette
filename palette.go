// Package palette fornece funcoes para visualizar paletas de cores no terminal.
// Utiliza lipgloss para renderizar swatches, grades ANSI 256, gradientes true color
// e temas personalizados lado a lado.
package palette

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

// ShowPalette renderiza swatches de cores com a largura especificada.
// Cada cor recebe um bloco colorido com seu valor hex abaixo.
func ShowPalette(colors []color.Color, width int) string {
	if len(colors) == 0 {
		return ""
	}

	swatchWidth := width / len(colors)
	if swatchWidth < 4 {
		swatchWidth = 4
	}

	var topRow, bottomRow strings.Builder

	for _, c := range colors {
		r, g, b, _ := c.RGBA()
		hex := fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)

		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(hex)).
			Width(swatchWidth).
			Render(strings.Repeat(" ", swatchWidth))

		label := lipgloss.NewStyle().
			Foreground(lipgloss.Color(hex)).
			Width(swatchWidth).
			Align(lipgloss.Center).
			Render(hex)

		topRow.WriteString(swatch)
		bottomRow.WriteString(label)
	}

	return topRow.String() + "\n" + bottomRow.String()
}

// ShowANSI256 exibe todas as 256 cores ANSI em uma grade organizada.
// As cores sao dispostas em grupos: cores padrao (0-15), cubo 6x6x6 (16-231)
// e escala de cinza (232-255).
func ShowANSI256() string {
	var sb strings.Builder

	sb.WriteString("  Standard Colors (0-15)\n")
	// Cores padrao 0-7
	for i := 0; i < 8; i++ {
		c := fmt.Sprintf("%d", i)
		style := lipgloss.NewStyle().
			Background(lipgloss.Color(c)).
			Width(4).
			Align(lipgloss.Center)
		sb.WriteString(style.Render(fmt.Sprintf("%3d", i)))
	}
	sb.WriteString("\n")
	// Cores claras 8-15
	for i := 8; i < 16; i++ {
		c := fmt.Sprintf("%d", i)
		style := lipgloss.NewStyle().
			Background(lipgloss.Color(c)).
			Foreground(lipgloss.Color("0")).
			Width(4).
			Align(lipgloss.Center)
		sb.WriteString(style.Render(fmt.Sprintf("%3d", i)))
	}
	sb.WriteString("\n\n")

	// Cubo de cores 16-231
	sb.WriteString("  Color Cube (16-231)\n")
	for row := 0; row < 12; row++ {
		for col := 0; col < 18; col++ {
			idx := 16 + row*18 + col
			if idx > 231 {
				break
			}
			c := fmt.Sprintf("%d", idx)
			fg := "15" // Branco como padrao
			if row > 5 {
				fg = "0" // Preto para cores claras
			}
			style := lipgloss.NewStyle().
				Background(lipgloss.Color(c)).
				Foreground(lipgloss.Color(fg)).
				Width(4).
				Align(lipgloss.Center)
			sb.WriteString(style.Render(fmt.Sprintf("%3d", idx)))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Escala de cinza 232-255
	sb.WriteString("  Grayscale (232-255)\n")
	for i := 232; i < 256; i++ {
		c := fmt.Sprintf("%d", i)
		fg := "15"
		if i > 243 {
			fg = "0"
		}
		style := lipgloss.NewStyle().
			Background(lipgloss.Color(c)).
			Foreground(lipgloss.Color(fg)).
			Width(4).
			Align(lipgloss.Center)
		sb.WriteString(style.Render(fmt.Sprintf("%3d", i)))
	}
	sb.WriteString("\n")

	return sb.String()
}

// ShowTrueColor renderiza um gradiente arco-iris usando cores true color.
// O gradiente varia a matiz horizontalmente e a luminosidade verticalmente.
func ShowTrueColor(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	var sb strings.Builder

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			hue := float64(x) / float64(width) * 360.0
			saturation := 1.0
			lightness := 0.2 + (float64(y)/float64(height))*0.6

			r, g, b := hslToRGB(hue, saturation, lightness)
			hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)

			style := lipgloss.NewStyle().
				Background(lipgloss.Color(hex))
			sb.WriteString(style.Render(" "))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ShowTheme renderiza uma paleta tematica com nome e cores nomeadas.
// Cada cor e exibida como um bloco com seu nome abaixo.
func ShowTheme(name string, colors map[string]color.Color) string {
	if len(colors) == 0 {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Padding(0, 1)

	var sb strings.Builder
	sb.WriteString(titleStyle.Render(fmt.Sprintf("Theme: %s", name)))
	sb.WriteString("\n\n")

	for label, c := range colors {
		r, g, b, _ := c.RGBA()
		hex := fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)

		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(hex)).
			Width(12).
			Render(strings.Repeat(" ", 12))

		labelStyled := lipgloss.NewStyle().
			Foreground(lipgloss.Color(hex)).
			Width(12).
			Render(label)

		hexStyled := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Width(12).
			Render(hex)

		sb.WriteString(fmt.Sprintf("  %s\n  %s\n  %s\n\n", swatch, labelStyled, hexStyled))
	}

	return sb.String()
}

// CompareThemes renderiza multiplos temas lado a lado para comparacao.
// Cada tema e exibido como uma coluna com suas cores empilhadas.
func CompareThemes(themes ...map[string]color.Color) string {
	if len(themes) == 0 {
		return ""
	}

	colWidth := 16
	var sb strings.Builder

	// Coleta todas as chaves unicas
	allKeys := make(map[string]bool)
	for _, theme := range themes {
		for k := range theme {
			allKeys[k] = true
		}
	}

	// Cabecalho com indice do tema
	for i := range themes {
		header := lipgloss.NewStyle().
			Bold(true).
			Width(colWidth).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("Theme %d", i+1))
		sb.WriteString(header)
		sb.WriteString("  ")
	}
	sb.WriteString("\n")

	// Para cada chave, renderiza uma linha com o swatch de cada tema
	for key := range allKeys {
		// Nome da cor
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render(fmt.Sprintf("%-12s", key)))
		sb.WriteString("\n")

		// Swatches
		for _, theme := range themes {
			c, ok := theme[key]
			if !ok {
				empty := lipgloss.NewStyle().
					Width(colWidth).
					Render(strings.Repeat("-", colWidth))
				sb.WriteString(empty)
			} else {
				r, g, b, _ := c.RGBA()
				hex := fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
				swatch := lipgloss.NewStyle().
					Background(lipgloss.Color(hex)).
					Width(colWidth).
					Render(strings.Repeat(" ", colWidth))
				sb.WriteString(swatch)
			}
			sb.WriteString("  ")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// hslToRGB converte valores HSL para RGB.
// h: 0-360, s: 0-1, l: 0-1
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}
