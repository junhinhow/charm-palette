# charm-palette

Visualizador e gerador de paletas de cores no terminal — pre-visualize temas diretamente no seu terminal.

> **[Read in English](README.md)**

## Funcionalidades

- **ShowPalette** — renderiza swatches de cores com valores hex
- **ShowANSI256** — exibe todas as 256 cores ANSI em grade
- **ShowTrueColor** — gradiente arco-iris com dimensoes configuraveis
- **ShowTheme** — exibe paleta tematica com cores nomeadas
- **CompareThemes** — comparacao de temas lado a lado

Construido com [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Instalacao

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

    // Gradiente arco-iris true color
    fmt.Println(palette.ShowTrueColor(80, 10))

    // Paleta personalizada
    cores := []color.Color{
        color.RGBA{R: 255, G: 89, B: 94, A: 255},
        color.RGBA{R: 255, G: 202, B: 58, A: 255},
        color.RGBA{R: 138, G: 201, B: 38, A: 255},
        color.RGBA{R: 25, G: 130, B: 196, A: 255},
    }
    fmt.Println(palette.ShowPalette(cores, 80))

    // Exibicao de tema
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
Renderiza swatches de cores em linha com labels hex abaixo.

### `ShowANSI256() string`
Exibe todas as 256 cores ANSI organizadas: padrao (0-15), cubo (16-231) e escala de cinza (232-255).

### `ShowTrueColor(width, height int) string`
Renderiza gradiente arco-iris com true color, variando matiz horizontal e luminosidade vertical.

### `ShowTheme(name string, colors map[string]color.Color) string`
Exibe um tema nomeado com swatches, labels e valores hex.

### `CompareThemes(themes ...map[string]color.Color) string`
Renderiza multiplos temas lado a lado para comparacao visual.

## Licenca

[MIT](LICENSE) - junhinhow
