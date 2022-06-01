// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 11:25
// version: 1.0.0
// desc   :

package ui

import (
	"fyne.io/fyne/v2"
	"gm3u8der/util"
)

const (
	iconPath = "./assets/img/icon.png"
)

func Icon(win fyne.Window) {
	icon := util.LoadResourceFromFile(iconPath)
	win.SetIcon(icon)
}
