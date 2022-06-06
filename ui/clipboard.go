// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-02 15:29
// version: 1.0.0
// desc   :

package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"golang.design/x/clipboard"
	"log"
	"regexp"
)

var (
	m3u8URLRegexp        = regexp.MustCompile("^https?://(.+)?(.m3u8)(\\?.*)?$")
	lastClipboardContent = ""
)

func Clipboard(win fyne.Window, showDialog func(m3u8URL string)) {
	err := clipboard.Init()
	if nil != err {
		log.Println(err)
		return
	}

	// 剪切板初始化成功，开始监听
	go watchClipboard(win, showDialog)
}

func watchClipboard(win fyne.Window, showDialog func(m3u8URL string)) {
	ch := clipboard.Watch(context.Background(), clipboard.FmtText)
	for data := range ch {
		content := string(data)
		if m3u8URLRegexp.MatchString(content) && content != lastClipboardContent {
			if nil != showDialog {
				showDialog(content)
			}
			lastClipboardContent = content
		}
	}
}
