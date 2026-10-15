// Package palette fornece funções para visualizar paletas de cores no terminal.
// Utiliza lipgloss para renderizar swatches, grades ANSI 256, gradientes true color
// e temas personalizados lado a lado.
package palette

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"reflect"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
)

// Erros retornados pelas variantes "Err" das funções públicas.
var (
	// ErrNoColors indica que nenhuma cor foi informada.
	ErrNoColors = errors.New("nenhuma cor informada")
	// ErrNoThemes indica que nenhum tema foi informado.
	ErrNoThemes = errors.New("nenhum tema informado")
	// ErrNilColor indica que a lista de cores contém uma entrada nula.
	ErrNilColor = errors.New("cor nula informada")
	// ErrInvalidWidth indica uma largura de exibição inválida.
	ErrInvalidWidth = errors.New("largura de exibição inválida")
)

const (
	// hexWidth é a largura de exibição de um rótulo hexadecimal no formato "#rrggbb".
	hexWidth = 7
	// minLabelWidth é a largura mínima da coluna de rótulos.
	minLabelWidth = 6
	// maxLabelWidth é a largura máxima da coluna de rótulos.
	maxLabelWidth = 24
	// themeSwatchWidth é a largura, em colunas, do swatch em ShowTheme.
	themeSwatchWidth = 12
	// compareColWidth é a largura, em colunas, de cada coluna de tema em CompareThemes.
	compareColWidth = 16
)

// amostra é uma cor pronta para renderização: a representação hexadecimal e uma
// cor de primeiro plano legível sobre ela.
type amostra struct {
	hex        string
	foreground color.Color
}

// ShowPalette renderiza swatches de cores com a largura especificada.
// Cada cor recebe um bloco colorido com seu valor hex abaixo.
//
// Entradas inválidas (cores nulas, largura não positiva) resultam em uma string
// vazia; use ShowPaletteErr para obter o detalhamento do erro.
func ShowPalette(colors []color.Color, width int) string {
	saida, err := ShowPaletteErr(colors, width)
	if err != nil {
		return ""
	}
	return saida
}

// ShowPaletteErr renderiza swatches de cores com a largura especificada e
// devolve um erro quando a entrada é inválida.
//
// O layout é horizontal quando há espaço para o rótulo "#rrggbb" de cada cor
// (pelo menos 7 colunas por cor). Caso contrário, as cores são empilhadas
// verticalmente, uma por linha, de modo que a largura solicitada nunca é
// ultrapassada.
func ShowPaletteErr(colors []color.Color, width int) (string, error) {
	if len(colors) == 0 {
		return "", ErrNoColors
	}
	if width <= 0 {
		return "", fmt.Errorf("%w: %d", ErrInvalidWidth, width)
	}

	amostras := make([]amostra, len(colors))
	for i, c := range colors {
		a, err := novaAmostra(c)
		if err != nil {
			return "", fmt.Errorf("cor no índice %d: %w", i, err)
		}
		amostras[i] = a
	}

	porAmostra := width / len(colors)
	if porAmostra < hexWidth {
		return renderPaletteVertical(amostras, width), nil
	}
	return renderPaletteHorizontal(amostras, porAmostra), nil
}

// renderPaletteHorizontal desenha todas as amostras em uma única linha de swatches
// com os rótulos hexadecimais centralizados na linha de baixo.
func renderPaletteHorizontal(amostras []amostra, largura int) string {
	var topo, base strings.Builder

	for _, a := range amostras {
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(a.hex)).
			Width(largura).
			Render(strings.Repeat(" ", largura))

		rotulo := lipgloss.NewStyle().
			Foreground(a.foreground).
			Width(largura).
			MaxWidth(largura).
			Align(lipgloss.Center).
			Render(a.hex)

		topo.WriteString(swatch)
		base.WriteString(rotulo)
	}

	return topo.String() + "\n" + base.String()
}

// renderPaletteVertical empilha as amostras, uma por linha, mantendo a largura
// solicitada. Usado quando não há espaço para o layout horizontal.
func renderPaletteVertical(amostras []amostra, largura int) string {
	linhas := make([]string, 0, len(amostras))

	for _, a := range amostras {
		if largura >= hexWidth+1 {
			larguraSwatch := largura - hexWidth - 1
			swatch := lipgloss.NewStyle().
				Background(lipgloss.Color(a.hex)).
				Width(larguraSwatch).
				Render(strings.Repeat(" ", larguraSwatch))
			linhas = append(linhas, swatch+" "+a.hex)
			continue
		}
		linhas = append(linhas, truncar(a.hex, largura))
	}

	return strings.Join(linhas, "\n")
}

// ShowANSI256 exibe todas as 256 cores ANSI em uma grade organizada.
// As cores são dispostas em grupos: cores padrão (0-15), cubo 6x6x6 (16-231)
// e escala de cinza (232-255).
func ShowANSI256() string {
	var sb strings.Builder

	sb.WriteString("  Standard Colors (0-15)\n")
	// Cores padrão 0-7
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
cubo:
	for row := 0; row < 12; row++ {
		for col := 0; col < 18; col++ {
			idx := 16 + row*18 + col
			if idx > 231 {
				// O break sem rótulo sairia apenas do laço interno,
				// continuando a renderizar as demais linhas do cubo.
				break cubo
			}
			c := fmt.Sprintf("%d", idx)
			fg := "15" // Branco como padrão
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

// ShowTrueColor renderiza um gradiente arco-íris usando cores true color.
// O gradiente varia a matiz horizontalmente e a luminosidade verticalmente.
//
// Dimensões inválidas resultam em uma string vazia; use ShowTrueColorErr para
// obter o detalhamento do erro.
func ShowTrueColor(width, height int) string {
	saida, err := ShowTrueColorErr(width, height)
	if err != nil {
		return ""
	}
	return saida
}

// ShowTrueColorErr renderiza um gradiente arco-íris e devolve um erro quando as
// dimensões informadas são inválidas.
func ShowTrueColorErr(width, height int) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("%w: %dx%d", ErrInvalidWidth, width, height)
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

	return sb.String(), nil
}

// ShowTheme renderiza uma paleta temática com nome e cores nomeadas.
// Cada cor é exibida como um bloco com seu nome e valor hexadecimal.
//
// As cores são ordenadas alfabeticamente, tornando a saída determinística.
// Entradas inválidas resultam em uma string vazia; use ShowThemeErr para obter
// o detalhamento do erro.
func ShowTheme(name string, colors map[string]color.Color) string {
	saida, err := ShowThemeErr(name, colors)
	if err != nil {
		return ""
	}
	return saida
}

// ShowThemeErr renderiza uma paleta temática e devolve um erro quando a entrada
// é inválida (lista vazia ou cor nula).
func ShowThemeErr(name string, colors map[string]color.Color) (string, error) {
	if len(colors) == 0 {
		return "", ErrNoColors
	}

	rotulos := chavesOrdenadas(colors)

	amostras := make([]amostra, len(rotulos))
	for i, rotulo := range rotulos {
		a, err := novaAmostra(colors[rotulo])
		if err != nil {
			return "", fmt.Errorf("cor %q: %w", rotulo, err)
		}
		amostras[i] = a
	}

	larguraRotulo := larguraRotulos(rotulos)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Padding(0, 1)

	var sb strings.Builder
	sb.WriteString(titleStyle.Render(fmt.Sprintf("Theme: %s", name)))
	sb.WriteString("\n\n")

	for i, rotulo := range rotulos {
		a := amostras[i]

		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(a.hex)).
			Width(themeSwatchWidth).
			Render(strings.Repeat(" ", themeSwatchWidth))

		// O rótulo usa contraste automático: desenhá-lo na própria cor
		// deixaria o texto ilegível em temas escuros.
		rotuloStyled := lipgloss.NewStyle().
			Foreground(a.foreground).
			Width(larguraRotulo).
			MaxWidth(larguraRotulo).
			Render(truncar(rotulo, larguraRotulo))

		hexStyled := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render(a.hex)

		sb.WriteString("  ")
		sb.WriteString(swatch)
		sb.WriteString("  ")
		sb.WriteString(rotuloStyled)
		sb.WriteString("  ")
		sb.WriteString(hexStyled)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// CompareThemes renderiza múltiplos temas lado a lado para comparação.
// Cada tema é exibido como uma coluna com suas cores empilhadas.
//
// As cores são ordenadas alfabeticamente, tornando a saída determinística.
// Entradas inválidas resultam em uma string vazia; use CompareThemesErr para
// obter o detalhamento do erro.
func CompareThemes(themes ...map[string]color.Color) string {
	saida, err := CompareThemesErr(themes...)
	if err != nil {
		return ""
	}
	return saida
}

// CompareThemesErr renderiza múltiplos temas lado a lado e devolve um erro
// quando algum tema contém uma cor nula.
func CompareThemesErr(themes ...map[string]color.Color) (string, error) {
	if len(themes) == 0 {
		return "", ErrNoThemes
	}

	// Coleta todas as chaves únicas
	todasChaves := make(map[string]bool)
	for i, theme := range themes {
		for chave, c := range theme {
			if ehNulo(c) {
				return "", fmt.Errorf("tema %d, cor %q: %w", i+1, chave, ErrNilColor)
			}
			todasChaves[chave] = true
		}
	}

	chaves := make([]string, 0, len(todasChaves))
	for chave := range todasChaves {
		chaves = append(chaves, chave)
	}
	sort.Strings(chaves)

	larguraRotulo := larguraRotulos(chaves)

	var sb strings.Builder

	// Cabeçalho com índice do tema, alinhado às colunas de cores.
	sb.WriteString(colunaRotulo("", larguraRotulo))
	for i := range themes {
		cabecalho := lipgloss.NewStyle().
			Bold(true).
			Width(compareColWidth).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("Theme %d", i+1))
		sb.WriteString(cabecalho)
		sb.WriteString("  ")
	}
	sb.WriteString("\n")

	// Para cada chave, renderiza uma linha com o swatch de cada tema
	for _, chave := range chaves {
		sb.WriteString(colunaRotulo(chave, larguraRotulo))
		for _, theme := range themes {
			c, ok := theme[chave]
			if !ok {
				vazio := lipgloss.NewStyle().
					Width(compareColWidth).
					Render(strings.Repeat("-", compareColWidth))
				sb.WriteString(vazio)
			} else {
				a, err := novaAmostra(c)
				if err != nil {
					return "", fmt.Errorf("cor %q: %w", chave, err)
				}
				swatch := lipgloss.NewStyle().
					Background(lipgloss.Color(a.hex)).
					Width(compareColWidth).
					Render(strings.Repeat(" ", compareColWidth))
				sb.WriteString(swatch)
			}
			sb.WriteString("  ")
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// colunaRotulo renderiza a célula da coluna de rótulos, com largura fixa e
// texto truncado para caber sem quebrar a linha.
func colunaRotulo(texto string, largura int) string {
	estilo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Width(largura).
		MaxWidth(largura)
	return estilo.Render(truncar(texto, largura)) + "  "
}

// chavesOrdenadas devolve as chaves de um mapa de cores em ordem alfabética,
// garantindo uma saída determinística.
func chavesOrdenadas(colors map[string]color.Color) []string {
	chaves := make([]string, 0, len(colors))
	for chave := range colors {
		chaves = append(chaves, chave)
	}
	sort.Strings(chaves)
	return chaves
}

// larguraRotulos calcula a largura da coluna de rótulos a partir do rótulo mais
// longo, respeitando os limites mínimo e máximo.
func larguraRotulos(rotulos []string) int {
	largura := 0
	for _, rotulo := range rotulos {
		if l := lipgloss.Width(rotulo); l > largura {
			largura = l
		}
	}
	if largura < minLabelWidth {
		largura = minLabelWidth
	}
	if largura > maxLabelWidth {
		largura = maxLabelWidth
	}
	return largura
}

// novaAmostra converte uma cor em sua representação hexadecimal e calcula a cor
// de primeiro plano legível sobre ela.
//
// Os canais são extraídos com o modelo NRGBA, ou seja, sem a multiplicação
// alfa que color.RGBA() aplica. Assim, cores sem opacidade total (NRGBA,
// RGBA64) não aparecem escurecidas.
func novaAmostra(c color.Color) (amostra, error) {
	if ehNulo(c) {
		return amostra{}, ErrNilColor
	}

	n := color.NRGBAModel.Convert(c).(color.NRGBA)

	return amostra{
		hex:        fmt.Sprintf("#%02x%02x%02x", n.R, n.G, n.B),
		foreground: contraste(n),
	}, nil
}

// contraste devolve a cor de primeiro plano (preto ou branco) com melhor
// contraste sobre a cor informada, segundo a luminância relativa da WCAG.
func contraste(n color.NRGBA) color.Color {
	// O ponto de equilíbrio entre texto preto e branco sobre um fundo de
	// luminância L é L ≈ 0.179.
	if luminancia(n) > 0.179 {
		return lipgloss.Color("#000000")
	}
	return lipgloss.Color("#ffffff")
}

// luminancia calcula a luminância relativa de uma cor conforme a WCAG.
func luminancia(n color.NRGBA) float64 {
	return 0.2126*canal(n.R) +
		0.7152*canal(n.G) +
		0.0722*canal(n.B)
}

// canal lineariza um canal de 8 bits conforme a WCAG.
func canal(v uint8) float64 {
	c := float64(v) / 255
	if c <= 0.03928 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

// ehNulo informa se a cor é nula, incluindo ponteiros nulos encapsulados em
// uma interface (por exemplo, (*color.RGBA)(nil)).
func ehNulo(c color.Color) bool {
	if c == nil {
		return true
	}
	v := reflect.ValueOf(c)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}

// truncar reduz o texto para no máximo max colunas de exibição, acrescentando
// reticências quando o conteúdo for cortado.
func truncar(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= max {
		return s
	}
	if max == 1 {
		return "…"
	}
	return lipgloss.NewStyle().MaxWidth(max-1).Render(s) + "…"
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
