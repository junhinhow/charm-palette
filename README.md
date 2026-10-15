# charm-palette

Terminal color palette viewer and generator — preview themes in your terminal.

> **[Leia em Português (PT-BR)](README.pt-br.md)**

## Features

- **ShowPalette** — render color swatches with hex labels
- **ShowANSI256** — display all 256 ANSI colors in a grid
- **ShowTrueColor** — rainbow gradient with configurable dimensions
- **ShowTheme** — themed palette display with named colors
- **CompareThemes** — side-by-side theme comparison

Powered by [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Every renderer is **deterministic**: themes are sorted alphabetically, so the
output never changes between calls (safe for snapshot tests). Colors are read
through the `NRGBA` model, so colors with partial alpha are not darkened, and
labels are painted with an **automatic contrast** color (black or white) instead
of the sample color itself, which would be unreadable on dark themes.

Invalid input (a `nil` color, a non-positive width) never panics: the `Show*`
functions return an empty string and the `Show*Err` variants return a detailed
error.

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
Renders color swatches in a row with hex labels below. When the requested width
is too small to fit the `#rrggbb` label of every swatch, the colors are stacked
vertically instead, so the requested width is never exceeded. Returns an empty
string for invalid input.

### `ShowPaletteErr(colors []color.Color, width int) (string, error)`
Same as `ShowPalette`, but reports invalid input (`ErrNoColors`, `ErrNilColor`,
`ErrInvalidWidth`).

### `ShowANSI256() string`
Displays all 256 ANSI colors organized by category: standard (0-15), color cube (16-231), and grayscale (232-255).

### `ShowTrueColor(width, height int) string`
Renders a rainbow gradient using true color, varying hue horizontally and lightness vertically.

### `ShowTrueColorErr(width, height int) (string, error)`
Same as `ShowTrueColor`, but reports invalid dimensions (`ErrInvalidWidth`).

### `ShowTheme(name string, colors map[string]color.Color) string`
Displays a named theme with color swatches, labels, and hex values. Colors are
sorted alphabetically. Returns an empty string for invalid input.

### `ShowThemeErr(name string, colors map[string]color.Color) (string, error)`
Same as `ShowTheme`, but reports invalid input (`ErrNoColors`, `ErrNilColor`).

### `CompareThemes(themes ...map[string]color.Color) string`
Renders multiple themes side by side for visual comparison, with the header
aligned to the color columns. Colors are sorted alphabetically. Returns an
empty string for invalid input.

### `CompareThemesErr(themes ...map[string]color.Color) (string, error)`
Same as `CompareThemes`, but reports invalid input (`ErrNoThemes`, `ErrNilColor`).

### Errors

| Error | Description |
|-------|-------------|
| `ErrNoColors` | No color was provided |
| `ErrNoThemes` | No theme was provided |
| `ErrNilColor` | The color list contains a `nil` entry |
| `ErrInvalidWidth` | The display width (or the gradient dimensions) is not positive |

## License

[MIT](LICENSE) - junhinhow
