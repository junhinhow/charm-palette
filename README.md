# charm-palette

Terminal color palette viewer and generator — preview themes in your terminal.

> **[Leia em Portugues (PT-BR)](README.pt-br.md)**

## Features

- **ShowPalette** — render color swatches with hex labels
- **ShowANSI256** — display all 256 ANSI colors in a grid
- **ShowTrueColor** — rainbow gradient with configurable dimensions
- **ShowTheme** — themed palette display with named colors
- **CompareThemes** — side-by-side theme comparison

Powered by [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Install

```bash
go get github.com/junhinhow/charm-palette@latest
```

## Usage

```go
package main

import (
    "fmt"
    "image/color"

    "github.com/junhinhow/charm-palette"
)

func main() {
    // Show ANSI 256 color grid
    fmt.Println(palette.ShowANSI256())

    // Show true color rainbow
    fmt.Println(palette.ShowTrueColor(80, 10))

    // Custom palette
    colors := []color.Color{
        color.RGBA{R: 255, G: 89, B: 94, A: 255},
        color.RGBA{R: 255, G: 202, B: 58, A: 255},
        color.RGBA{R: 138, G: 201, B: 38, A: 255},
        color.RGBA{R: 25, G: 130, B: 196, A: 255},
        color.RGBA{R: 106, G: 76, B: 147, A: 255},
    }
    fmt.Println(palette.ShowPalette(colors, 80))

    // Theme display
    theme := map[string]color.Color{
        "primary":    color.RGBA{R: 99, G: 102, B: 241, A: 255},
        "secondary":  color.RGBA{R: 168, G: 85, B: 247, A: 255},
        "background": color.RGBA{R: 15, G: 23, B: 42, A: 255},
        "text":       color.RGBA{R: 226, G: 232, B: 240, A: 255},
    }
    fmt.Println(palette.ShowTheme("Indigo Night", theme))
}
```

## API

### `ShowPalette(colors []color.Color, width int) string`
Renders color swatches in a row with hex labels below.

### `ShowANSI256() string`
Displays all 256 ANSI colors organized by category: standard (0-15), color cube (16-231), and grayscale (232-255).

### `ShowTrueColor(width, height int) string`
Renders a rainbow gradient using true color, varying hue horizontally and lightness vertically.

### `ShowTheme(name string, colors map[string]color.Color) string`
Displays a named theme with color swatches, labels, and hex values.

### `CompareThemes(themes ...map[string]color.Color) string`
Renders multiple themes side by side for visual comparison.

## License

[MIT](LICENSE) - junhinhow
