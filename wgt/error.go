// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-02 9:29
// version: 1.0.0
// desc   :

package wgt

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func ShowErrorDialog(title, message string, win fyne.Window, okCallback func()) {
	dialog.NewConfirm(title, message, func(ok bool) {
		if ok {
			okCallback()
		}
	}, win).Show()
}
