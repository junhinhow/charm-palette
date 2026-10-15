# charm-palette

Visualizador e gerador de paletas de cores no terminal — pre-visualize temas diretamente no seu terminal.

> **[Read in English](README.md)**

## Funcionalidades

- **ShowPalette** — renderiza swatches de cores com valores hex
- **ShowANSI256** — exibe todas as 256 cores ANSI em grade
- **ShowTrueColor** — gradiente arco-íris com dimensões configuráveis
- **ShowTheme** — exibe paleta temática com cores nomeadas
- **CompareThemes** — comparação de temas lado a lado

Construído com [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Todos os renderizadores são **determinísticos**: os temas são ordenados
alfabeticamente, de modo que a saída nunca muda entre chamadas (seguro para
testes por snapshot). As cores são lidas através do modelo `NRGBA`, portanto
cores com alfa parcial não aparecem escurecidas, e os rótulos são desenhados
com uma cor de **contraste automático** (preto ou branco) em vez da própria cor
da amostra, que seria ilegível em temas escuros.

Entradas inválidas (uma cor `nil`, uma largura não positiva) nunca causam
pânico: as funções `Show*` devolvem uma string vazia e as variantes `Show*Err`
devolvem um erro detalhado.

## Instalação

```bash
go get github.com/junhinhow/charm-palette@latest
```

## Uso

```go
package main

import (
    "fmt"
    "image/color"

    "github.com/junhinhow/charm-palette"
)

func main() {
    // Grade de cores ANSI 256
    fmt.Println(palette.ShowANSI256())

    // Gradiente arco-íris true color
    fmt.Println(palette.ShowTrueColor(80, 10))

    // Paleta personalizada
    cores := []color.Color{
        color.RGBA{R: 255, G: 89, B: 94, A: 255},
        color.RGBA{R: 255, G: 202, B: 58, A: 255},
        color.RGBA{R: 138, G: 201, B: 38, A: 255},
        color.RGBA{R: 25, G: 130, B: 196, A: 255},
    }
    fmt.Println(palette.ShowPalette(cores, 80))

    // Exibição de tema
    tema := map[string]color.Color{
        "primaria":  color.RGBA{R: 99, G: 102, B: 241, A: 255},
        "secundaria": color.RGBA{R: 168, G: 85, B: 247, A: 255},
        "fundo":     color.RGBA{R: 15, G: 23, B: 42, A: 255},
        "texto":     color.RGBA{R: 226, G: 232, B: 240, A: 255},
    }
    fmt.Println(palette.ShowTheme("Noite Indigo", tema))
}
```

## API

### `ShowPalette(colors []color.Color, width int) string`
Renderiza swatches de cores em linha com rótulos hex abaixo. Quando a largura
solicitada é pequena demais para o rótulo `#rrggbb` de cada cor, as cores são
empilhadas verticalmente, de modo que a largura pedida nunca é ultrapassada.
Devolve uma string vazia para entradas inválidas.

### `ShowPaletteErr(colors []color.Color, width int) (string, error)`
Igual a `ShowPalette`, mas informa as entradas inválidas (`ErrNoColors`,
`ErrNilColor`, `ErrInvalidWidth`).

### `ShowANSI256() string`
Exibe todas as 256 cores ANSI organizadas: padrão (0-15), cubo (16-231) e escala de cinza (232-255).

### `ShowTrueColor(width, height int) string`
Renderiza gradiente arco-íris com true color, variando matiz horizontal e luminosidade vertical.

### `ShowTrueColorErr(width, height int) (string, error)`
Igual a `ShowTrueColor`, mas informa dimensões inválidas (`ErrInvalidWidth`).

### `ShowTheme(name string, colors map[string]color.Color) string`
Exibe um tema nomeado com swatches, rótulos e valores hex. As cores são
ordenadas alfabeticamente. Devolve uma string vazia para entradas inválidas.

### `ShowThemeErr(name string, colors map[string]color.Color) (string, error)`
Igual a `ShowTheme`, mas informa as entradas inválidas (`ErrNoColors`,
`ErrNilColor`).

### `CompareThemes(themes ...map[string]color.Color) string`
Renderiza múltiplos temas lado a lado para comparação visual, com o cabeçalho
alinhado às colunas de cores. As cores são ordenadas alfabeticamente. Devolve
uma string vazia para entradas inválidas.

### `CompareThemesErr(themes ...map[string]color.Color) (string, error)`
Igual a `CompareThemes`, mas informa as entradas inválidas (`ErrNoThemes`,
`ErrNilColor`).

### Erros

| Erro | Descrição |
|------|-----------|
| `ErrNoColors` | Nenhuma cor foi informada |
| `ErrNoThemes` | Nenhum tema foi informado |
| `ErrNilColor` | A lista de cores contém uma entrada `nil` |
| `ErrInvalidWidth` | A largura de exibição (ou as dimensões do gradiente) não é positiva |

## Licença

[MIT](LICENSE) - junhinhow
