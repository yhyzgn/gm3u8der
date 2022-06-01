// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-05-31 12:48
// version: 1.0.0
// desc   :

package thm

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"strings"
)

type FontTheme struct {
	regular, bold, italic, boldItalic, monospace fyne.Resource
}

func (ft *FontTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground && variant == theme.VariantLight {
		return color.White
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (ft *FontTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (ft *FontTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return ft.monospace
	}
	if style.Bold {
		if style.Italic {
			return ft.boldItalic
		}
		return ft.bold
	}
	if style.Italic {
		return ft.italic
	}
	return ft.regular
}

func (ft *FontTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (ft *FontTheme) SetFonts(regularFontPath string, monoFontPath string) {
	if regularFontPath != "" {
		ft.regular = loadFont(regularFontPath, "Regular", ft.regular)
		ft.bold = loadFont(regularFontPath, "Bold", ft.bold)
		ft.italic = loadFont(regularFontPath, "Italic", ft.italic)
		ft.boldItalic = loadFont(regularFontPath, "Bold Italic", ft.boldItalic)
	}
	if monoFontPath != "" {
		ft.monospace = loadFont(monoFontPath, "Regular", ft.boldItalic)
	}
}

func NewFontTheme() *FontTheme {
	ft := &FontTheme{
		regular:    theme.TextFont(),
		bold:       theme.TextBoldFont(),
		italic:     theme.TextItalicFont(),
		boldItalic: theme.TextBoldItalicFont(),
		monospace:  theme.TextMonospaceFont(),
	}
	ft.SetFonts("./assets/ttf/Consolas-with-Yahei Regular Nerd Font.ttf", "")
	return ft
}

func loadFont(env, variant string, fallback fyne.Resource) fyne.Resource {
	variantPath := strings.ReplaceAll(env, "Regular", variant)
	res, err := fyne.LoadResourceFromPath(variantPath)
	if nil != err {
		fyne.LogError("Error loading specified font", err)
		return fallback
	}
	return res
}
