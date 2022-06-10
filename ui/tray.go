// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-10 9:42
// version: 1.0.0
// desc   :

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"gm3u8der/holder"
)

var winShowing = true

// SetTray ...
func SetTray(app fyne.App, win fyne.Window) {
	win.SetCloseIntercept(func() {
		if holder.Settings.TraySticky {
			win.Hide()
			winShowing = false
		} else {
			win.Close()
		}
	})

	mi := fyne.NewMenuItem(showMainWindowLabel(), func() {
		if nil != win {
			if winShowing {
				win.Hide()
				winShowing = false
			} else {
				win.Show()
				winShowing = true
			}
		}
	})

	if desk, ok := app.(desktop.App); ok {
		menu := fyne.NewMenu(showMainWindowLabel(), mi)
		desk.SetSystemTrayMenu(menu)
	}
}

func showMainWindowLabel() string {
	return "显示/隐藏主界面"
}
