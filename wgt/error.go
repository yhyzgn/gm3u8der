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

func ShowErrorDialog(title, message string, win fyne.Window, okCallback func(), cancelCallback func()) {
	dialog.NewConfirm(title, message, func(innerOk bool) {
		if innerOk {
			if nil != okCallback {
				okCallback()
			}
		} else {
			if nil != cancelCallback {
				cancelCallback()
			}
		}
	}, win).Show()
}
