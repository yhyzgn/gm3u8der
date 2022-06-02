// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-02 15:29
// version: 1.0.0
// desc   :

package ui

import (
	"fyne.io/fyne/v2"
	"gm3u8der/component"
	"regexp"
	"time"
)

var (
	m3u8URLRegexp        = regexp.MustCompile("^https?://(.+)?(.m3u8)(\\?.*)?$")
	lastClipboardContent = ""
)

func Clipboard(win fyne.Window, showDialog func(m3u8URL string)) {
	go watchClipboard(win, showDialog)
}

func watchClipboard(win fyne.Window, showDialog func(m3u8URL string)) {
	component.StartTicker(time.Second, func() {
		content := win.Clipboard().Content()
		if m3u8URLRegexp.MatchString(content) && content != lastClipboardContent {
			if nil != showDialog {
				showDialog(content)
			}
			lastClipboardContent = content
		}
	})
}
