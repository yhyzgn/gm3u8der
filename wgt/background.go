// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-05-31 13:56
// version: 1.0.0
// desc   :

package wgt

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"image/color"
)

// SetBackgroundColor 给控件设置背景颜色
func SetBackgroundColor(co fyne.CanvasObject, clr color.Color) fyne.CanvasObject {
	return container.New(layout.NewMaxLayout(), canvas.NewRectangle(clr), co)
}
