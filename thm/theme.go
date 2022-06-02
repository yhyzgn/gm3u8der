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
)

type FontTheme struct{}

var _ fyne.Theme = (*FontTheme)(nil)

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
		return theme.DefaultTheme().Font(style)
	}
	if style.Bold {
		if style.Italic {
			return theme.DefaultTheme().Font(style)
		}
		return resourceLXGWWenKaiRegularTtf
	}
	if style.Italic {
		return theme.DefaultTheme().Font(style)
	}
	return resourceLXGWWenKaiRegularTtf
}

func (ft *FontTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func NewFontTheme() *FontTheme {
	return &FontTheme{}
}
