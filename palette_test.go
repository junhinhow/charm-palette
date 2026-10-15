package palette

import (
	"errors"
	"fmt"
	"image/color"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// coresDeTeste devolve uma paleta varietyada usada pelos testes.
func coresDeTeste() []color.Color {
	return []color.Color{
		color.RGBA{R: 255, G: 89, B: 94, A: 255},
		color.RGBA{R: 255, G: 202, B: 58, A: 255},
		color.RGBA{R: 138, G: 201, B: 38, A: 255},
		color.RGBA{R: 25, G: 130, B: 196, A: 255},
		color.RGBA{R: 106, G: 76, B: 147, A: 255},
		color.RGBA{R: 1, G: 2, B: 3, A: 255},
	}
}

// colunaDe devolve a coluna de exibição de uma substring dentro de um texto
// estilizado (códigos ANSI ignorados). Devolve -1 quando a substring não existe.
func colunaDe(texto, busca string) int {
	idx := strings.Index(texto, busca)
	if idx < 0 {
		return -1
	}
	return lipgloss.Width(texto[:idx])
}

// colunaDoPrimeiroFundo devolve a coluna de exibição em que começa o primeiro
// bloco com cor de fundo de uma linha estilizada.
func colunaDoPrimeiroFundo(linha string) int {
	idx := strings.Index(linha, "\x1b[48;2;")
	if idx < 0 {
		return -1
	}
	return lipgloss.Width(linha[:idx])
}

// TestShowPaletteCorNulaPanica garante que uma cor nula não causa pânico e que o
// erro devolvido identifica o problema.
func TestShowPaletteCorNulaPanica(t *testing.T) {
	saida := ShowPalette([]color.Color{nil}, 20)
	if saida != "" {
		t.Fatalf("esperava string vazia para cor nula, obtive %q", saida)
	}

	_, err := ShowPaletteErr([]color.Color{nil}, 20)
	if !errors.Is(err, ErrNilColor) {
		t.Fatalf("esperava ErrNilColor, obtive %v", err)
	}
}

// TestShowPalettePonteiroNuloPanica garante que um ponteiro nil encapsulado em
// uma interface também é tratado como cor nula.
func TestShowPalettePonteiroNuloPanica(t *testing.T) {
	var nulo *color.RGBA

	if saida := ShowPalette([]color.Color{nulo}, 20); saida != "" {
		t.Fatalf("esperava string vazia para ponteiro nulo, obtive %q", saida)
	}

	if _, err := ShowPaletteErr([]color.Color{nulo}, 20); !errors.Is(err, ErrNilColor) {
		t.Fatalf("esperava ErrNilColor, obtive %v", err)
	}
}

// TestShowThemeCorNulaPanica garante que ShowTheme e CompareThemes não aceitam
// cores nulas.
func TestShowThemeCorNulaPanica(t *testing.T) {
	tema := map[string]color.Color{"primaria": color.RGBA{R: 1, G: 2, B: 3, A: 255}, "nula": nil}

	if saida := ShowTheme("Teste", tema); saida != "" {
		t.Fatalf("esperava string vazia para cor nula, obtive %q", saida)
	}
	if _, err := ShowThemeErr("Teste", tema); !errors.Is(err, ErrNilColor) {
		t.Fatalf("esperava ErrNilColor, obtive %v", err)
	}

	if saida := CompareThemes(tema); saida != "" {
		t.Fatalf("esperava string vazia para cor nula, obtive %q", saida)
	}
	if _, err := CompareThemesErr(tema); !errors.Is(err, ErrNilColor) {
		t.Fatalf("esperava ErrNilColor, obtive %v", err)
	}
}

// TestShowPaletteLarguraPequena garante que a largura solicitada é respeitada
// e que o rótulo hexadecimal nunca é quebrado em duas linhas.
func TestShowPaletteLarguraPequena(t *testing.T) {
	cores := coresDeTeste()
	const largura = 20

	saida := ShowPalette(cores, largura)
	if saida == "" {
		t.Fatal("esperava saída não vazia")
	}

	for _, linha := range strings.Split(saida, "\n") {
		if w := lipgloss.Width(linha); w > largura {
			t.Fatalf("linha com %d colunas ultrapassa a largura %d: %q", w, largura, linha)
		}
	}

	for _, c := range cores {
		n := color.NRGBAModel.Convert(c).(color.NRGBA)
		hex := fmt.Sprintf("#%02x%02x%02x", n.R, n.G, n.B)
		if !strings.Contains(saida, hex) {
			t.Fatalf("rótulo %q não encontrado na saída (quebrado?): %q", hex, saida)
		}
	}
}

// TestShowPaletteLarguraHorizontalExata garante que, quando há espaço, os
// swatches ocupam exatamente a largura solicitada.
func TestShowPaletteLarguraHorizontalExata(t *testing.T) {
	cores := coresDeTeste()[:5]

	saida := ShowPalette(cores, 80)
	linhas := strings.Split(saida, "\n")
	if len(linhas) != 2 {
		t.Fatalf("esperava 2 linhas no layout horizontal, obtive %d", len(linhas))
	}
	for _, linha := range linhas {
		if w := lipgloss.Width(linha); w != 80 {
			t.Fatalf("esperava 80 colunas, obtive %d: %q", w, linha)
		}
	}
}

// TestShowPaletteLarguraInvalida garante que larguras não positivas geram erro.
func TestShowPaletteLarguraInvalida(t *testing.T) {
	if _, err := ShowPaletteErr(coresDeTeste(), 0); !errors.Is(err, ErrInvalidWidth) {
		t.Fatalf("esperava ErrInvalidWidth, obtive %v", err)
	}
	if saida := ShowPalette(coresDeTeste(), -5); saida != "" {
		t.Fatalf("esperava string vazia para largura inválida, obtive %q", saida)
	}
}

// TestShowPaletteAlfaNaoPreMultiplicado garante que os canais usados são os
// diretos (sem multiplicação alfa), evitando cores escurecidas.
func TestShowPaletteAlfaNaoPreMultiplicado(t *testing.T) {
	casos := []struct {
		nome string
		cor  color.Color
		hex  string
	}{
		{"NRGBA sem opacidade", color.NRGBA{R: 255, G: 0, B: 0, A: 128}, "#ff0000"},
		{"NRGBA opaco", color.NRGBA{R: 18, G: 52, B: 86, A: 255}, "#123456"},
		{"RGBA64", color.RGBA64{R: 0x1234, G: 0x5678, B: 0x9abc, A: 0xffff}, "#12569a"},
	}

	for _, c := range casos {
		saida := ShowPalette([]color.Color{c.cor}, 20)
		if !strings.Contains(saida, c.hex) {
			t.Fatalf("%s: esperava %q na saída, obtive %q", c.nome, c.hex, saida)
		}
	}
}

// TestShowThemeOrdemDeterministica garante que a ordem de renderização das
// cores não muda entre chamadas (a iteração de map é aleatória).
func TestShowThemeOrdemDeterministica(t *testing.T) {
	tema := map[string]color.Color{
		"background": color.RGBA{R: 15, G: 23, B: 42, A: 255},
		"primary":    color.RGBA{R: 99, G: 102, B: 241, A: 255},
		"secondary":  color.RGBA{R: 168, G: 85, B: 247, A: 255},
		"text":       color.RGBA{R: 226, G: 232, B: 240, A: 255},
		"border":     color.RGBA{R: 51, G: 65, B: 85, A: 255},
		"muted":      color.RGBA{R: 100, G: 116, B: 139, A: 255},
	}

	esperado := ShowTheme("Ordem", tema)
	if esperado == "" {
		t.Fatal("esperava saída não vazia")
	}

	for i := 0; i < 50; i++ {
		if obtida := ShowTheme("Ordem", tema); obtida != esperado {
			t.Fatalf("saída mudou entre chamadas (iteração %d):\n%q\n%q", i, esperado, obtida)
		}
	}

	// As cores devem aparecer em ordem alfabética.
	posBackground := strings.Index(esperado, "background")
	posPrimary := strings.Index(esperado, "primary")
	posSecondary := strings.Index(esperado, "secondary")
	if posBackground < 0 || posPrimary < 0 || posSecondary < 0 {
		t.Fatalf("esperava encontrar os rótulos na saída: %q", esperado)
	}
	if !(posBackground < posPrimary && posPrimary < posSecondary) {
		t.Fatalf("esperava ordem alfabética dos rótulos: %d, %d, %d", posBackground, posPrimary, posSecondary)
	}
}

// TestCompareThemesOrdemDeterministica garante que a ordem das linhas de
// CompareThemes é estável entre chamadas.
func TestCompareThemesOrdemDeterministica(t *testing.T) {
	temaA := map[string]color.Color{
		"background": color.RGBA{R: 15, G: 23, B: 42, A: 255},
		"primary":    color.RGBA{R: 99, G: 102, B: 241, A: 255},
		"secondary":  color.RGBA{R: 168, G: 85, B: 247, A: 255},
		"text":       color.RGBA{R: 226, G: 232, B: 240, A: 255},
	}
	temaB := map[string]color.Color{
		"background": color.RGBA{R: 0, G: 0, B: 0, A: 255},
		"primary":    color.RGBA{R: 220, G: 38, B: 38, A: 255},
		"secondary":  color.RGBA{R: 34, G: 197, B: 94, A: 255},
		"text":       color.RGBA{R: 250, G: 204, B: 21, A: 255},
	}

	esperado := CompareThemes(temaA, temaB)
	if esperado == "" {
		t.Fatal("esperava saída não vazia")
	}

	for i := 0; i < 50; i++ {
		if obtida := CompareThemes(temaA, temaB); obtida != esperado {
			t.Fatalf("saída mudou entre chamadas (iteração %d):\n%q\n%q", i, esperado, obtida)
		}
	}
}

// TestCompareThemesCabecalhoAlinhado garante que o cabeçalho de tema é
// alinhado com as colunas de cores.
func TestCompareThemesCabecalhoAlinhado(t *testing.T) {
	tema := map[string]color.Color{
		"background": color.RGBA{R: 15, G: 23, B: 42, A: 255},
		"primary":    color.RGBA{R: 99, G: 102, B: 241, A: 255},
	}
	segundo := map[string]color.Color{
		"background": color.RGBA{R: 0, G: 0, B: 0, A: 255},
		"primary":    color.RGBA{R: 220, G: 38, B: 38, A: 255},
	}

	saida := CompareThemes(tema, segundo)
	linhas := strings.Split(strings.TrimRight(saida, "\n"), "\n")
	if len(linhas) < 2 {
		t.Fatalf("esperava cabeçalho e linhas de dados, obtive %d linhas", len(linhas))
	}

	cabecalho := linhas[0]
	colunaCabecalho := colunaDe(cabecalho, "Theme 1")
	if colunaCabecalho < 0 {
		t.Fatalf("cabeçalho sem o rótulo do tema: %q", cabecalho)
	}

	// Primeira linha que contém um swatch (primeira linha de dados).
	var dados string
	for _, linha := range linhas[1:] {
		if strings.Contains(linha, "\x1b[48;2;") {
			dados = linha
			break
		}
	}
	if dados == "" {
		t.Fatalf("nenhuma linha de dados com swatch: %q", saida)
	}

	colunaSwatch := colunaDoPrimeiroFundo(dados)
	if colunaSwatch < 0 {
		t.Fatalf("não foi possível localizar o primeiro swatch: %q", dados)
	}

	// O rótulo do tema deve ficar dentro da coluna do swatch.
	if colunaCabecalho < colunaSwatch {
		t.Fatalf("cabeçalho começa na coluna %d, antes do swatch na coluna %d: %q",
			colunaCabecalho, colunaSwatch, cabecalho)
	}
	if colunaCabecalho+len("Theme 1") > colunaSwatch+compareColWidth {
		t.Fatalf("cabeçalho termina na coluna %d, fora da coluna do swatch (%d..%d): %q",
			colunaCabecalho+len("Theme 1"), colunaSwatch, colunaSwatch+compareColWidth, cabecalho)
	}

	// Todas as linhas devem ter a mesma largura para que as colunas alinhem.
	largura := lipgloss.Width(linhas[0])
	for _, linha := range linhas {
		if w := lipgloss.Width(linha); w != largura {
			t.Fatalf("linha com %d colunas difere do cabeçalho (%d): %q", w, largura, linha)
		}
	}
}

// TestShowThemeRotuloContrasteAutomatico garante que o rótulo não é desenhado
// na própria cor da amostra (ilegível em cores escuras).
func TestShowThemeRotuloContrasteAutomatico(t *testing.T) {
	casos := []struct {
		nome     string
		cor      color.Color
		proibido string
		esperado string
	}{
		{
			nome:     "cor escura",
			cor:      color.RGBA{R: 15, G: 23, B: 42, A: 255},
			proibido: "\x1b[38;2;15;23;42m",
			esperado: "\x1b[38;2;255;255;255m",
		},
		{
			nome:     "cor clara",
			cor:      color.RGBA{R: 240, G: 230, B: 140, A: 255},
			proibido: "\x1b[38;2;240;230;140m",
			esperado: "\x1b[38;2;0;0;0m",
		},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			saida := ShowTheme("Contraste", map[string]color.Color{"amostra": c.cor})

			if strings.Contains(saida, c.proibido) {
				t.Fatalf("rótulo desenhado na própria cor da amostra: %q", saida)
			}
			if !strings.Contains(saida, c.esperado) {
				t.Fatalf("rótulo não usa o contraste automático (%q): %q", c.esperado, saida)
			}
		})
	}
}

// TestShowPaletteRotuloContrasteAutomatico garante a mesma garantia no
// ShowPalette, onde o rótulo hex é desenhado embaixo do swatch.
func TestShowPaletteRotuloContrasteAutomatico(t *testing.T) {
	saida := ShowPalette([]color.Color{color.RGBA{R: 15, G: 23, B: 42, A: 255}}, 20)
	if strings.Contains(saida, "\x1b[38;2;15;23;42m") {
		t.Fatalf("rótulo desenhado na própria cor da amostra: %q", saida)
	}
	if !strings.Contains(saida, "\x1b[38;2;255;255;255m") {
		t.Fatalf("rótulo não usa o contraste automático: %q", saida)
	}
}

// TestRotulosLongosNaoQuebramLinha garante que rótulos com acentos ou muito
// longos não quebram o layout (o padding deve considerar a largura de
// exibição, não a contagem de bytes).
func TestRotulosLongosNaoQuebramLinha(t *testing.T) {
	tema := map[string]color.Color{
		"configuração":                 color.RGBA{R: 99, G: 102, B: 241, A: 255},
		"nome-de-etapa-muito-comprido": color.RGBA{R: 168, G: 85, B: 247, A: 255},
	}

	saida := ShowTheme("Longos", tema)
	for _, linha := range strings.Split(saida, "\n") {
		if strings.Contains(linha, "\n") {
			t.Fatalf("linha com quebra inesperada: %q", linha)
		}
	}

	// Nenhuma linha pode conter duas amostras (o swatch e o rótulo).
	for _, linha := range strings.Split(saida, "\n") {
		if n := strings.Count(linha, "\x1b[48;2;"); n > 1 {
			t.Fatalf("linha com %d swatches, layout quebrado: %q", n, linha)
		}
	}

	// O rótulo longo deve ser truncado, não estourar a coluna.
	if !strings.Contains(saida, "…") {
		t.Fatalf("esperava rótulo truncado com reticências: %q", saida)
	}
}

// TestCompareThemesRotulosLongosNaoQuebramLinha garante o mesmo para o
// CompareThemes, inclusive com rótulos multibyte.
func TestCompareThemesRotulosLongosNaoQuebramLinha(t *testing.T) {
	tema := map[string]color.Color{
		"configuração":                 color.RGBA{R: 99, G: 102, B: 241, A: 255},
		"nome-de-etapa-muito-comprido": color.RGBA{R: 168, G: 85, B: 247, A: 255},
	}

	saida := CompareThemes(tema)
	linhas := strings.Split(strings.TrimRight(saida, "\n"), "\n")

	largura := lipgloss.Width(linhas[0])
	for _, linha := range linhas {
		if w := lipgloss.Width(linha); w != largura {
			t.Fatalf("linha com %d colunas difere do cabeçalho (%d): %q", w, largura, linha)
		}
	}
}

// TestHSLToRGB cobre os seis setores da roda de cores e o limite h=360.
func TestHSLToRGB(t *testing.T) {
	casos := []struct {
		h       float64
		r, g, b uint8
	}{
		{0, 255, 0, 0},
		{60, 255, 255, 0},
		{120, 0, 255, 0},
		{180, 0, 255, 255},
		{240, 0, 0, 255},
		{300, 255, 0, 255},
		{360, 255, 0, 0},
	}

	for _, c := range casos {
		r, g, b := hslToRGB(c.h, 1, 0.5)
		if r != c.r || g != c.g || b != c.b {
			t.Errorf("hslToRGB(%v, 1, 0.5) = (%d, %d, %d), esperado (%d, %d, %d)",
				c.h, r, g, b, c.r, c.g, c.b)
		}
	}

	// h=360 deve equivaler a h=0.
	r360, g360, b360 := hslToRGB(360, 1, 0.5)
	r0, g0, b0 := hslToRGB(0, 1, 0.5)
	if r360 != r0 || g360 != g0 || b360 != b0 {
		t.Errorf("h=360 = (%d, %d, %d), esperado igual a h=0 (%d, %d, %d)",
			r360, g360, b360, r0, g0, b0)
	}

	// Saturação zero devolve cinza.
	r, g, b := hslToRGB(123, 0, 0.5)
	if r != 127 && r != 128 {
		t.Errorf("hslToRGB com saturação 0 esperava cinza, obtive (%d, %d, %d)", r, g, b)
	}
	if r != g || g != b {
		t.Errorf("hslToRGB com saturação 0 esperava canais iguais, obtive (%d, %d, %d)", r, g, b)
	}
}

// TestShowANSI256Grade garante que a grade ANSI 256 tem a estrutura esperada:
// 8 cores padrão, 8 claras, 12 linhas de 18 no cubo e 24 na escala de cinza.
func TestShowANSI256Grade(t *testing.T) {
	saida := ShowANSI256()
	linhas := strings.Split(strings.TrimRight(saida, "\n"), "\n")

	inicioCubo := -1
	fimCubo := -1
	for i, linha := range linhas {
		if inicioCubo < 0 && strings.Contains(linha, "Color Cube") {
			inicioCubo = i
			continue
		}
		if inicioCubo >= 0 && strings.Contains(linha, "Grayscale") {
			fimCubo = i
			break
		}
	}

	if inicioCubo < 0 || fimCubo < 0 {
		t.Fatalf("não foi possível localizar as seções da grade: %q", saida)
	}

	// Entre o cabeçalho do cubo e o cabeçalho da escala de cinza há uma linha
	// em branco e 12 linhas com 18 células de 4 colunas.
	linhasCubo := linhas[inicioCubo+1 : fimCubo-1]
	if len(linhasCubo) != 12 {
		t.Fatalf("esperava 12 linhas no cubo de cores, obtive %d", len(linhasCubo))
	}
	for i, linha := range linhasCubo {
		if w := lipgloss.Width(linha); w != 18*4 {
			t.Fatalf("linha %d do cubo com %d colunas, esperado %d", i, w, 18*4)
		}
	}
}

// TestShowTrueColorDimensoes garante que dimensões inválidas geram erro.
func TestShowTrueColorDimensoes(t *testing.T) {
	if saida := ShowTrueColor(0, 10); saida != "" {
		t.Fatalf("esperava string vazia para largura inválida, obtive %q", saida)
	}
	if _, err := ShowTrueColorErr(10, 0); !errors.Is(err, ErrInvalidWidth) {
		t.Fatalf("esperava ErrInvalidWidth, obtive %v", err)
	}

	saida := ShowTrueColor(4, 2)
	linhas := strings.Split(strings.TrimRight(saida, "\n"), "\n")
	if len(linhas) != 2 {
		t.Fatalf("esperava 2 linhas no gradiente, obtive %d", len(linhas))
	}
	for _, linha := range linhas {
		if w := lipgloss.Width(linha); w != 4 {
			t.Fatalf("linha do gradiente com %d colunas, esperado 4", w)
		}
	}
}

// TestShowThemeSemCores e TestCompareThemesSemTemas garantem o tratamento das
// entradas vazias.
func TestShowThemeSemCores(t *testing.T) {
	if saida := ShowTheme("Vazio", nil); saida != "" {
		t.Fatalf("esperava string vazia, obtive %q", saida)
	}
	if _, err := ShowThemeErr("Vazio", nil); !errors.Is(err, ErrNoColors) {
		t.Fatalf("esperava ErrNoColors, obtive %v", err)
	}
	if saida := CompareThemes(); saida != "" {
		t.Fatalf("esperava string vazia, obtive %q", saida)
	}
	if _, err := CompareThemesErr(); !errors.Is(err, ErrNoThemes) {
		t.Fatalf("esperava ErrNoThemes, obtive %v", err)
	}
}
