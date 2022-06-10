// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 11:25
// version: 1.0.0
// desc   :

package ui

import (
	"fyne.io/fyne/v2"
)

// Icon ...
func Icon(app fyne.App, win fyne.Window) {
	app.SetIcon(resourceLogoPng)
	win.SetIcon(resourceLogoPng)
}
