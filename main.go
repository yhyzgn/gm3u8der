// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-05-30 15:35
// version: 1.0.0
// desc   :

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"gm3u8der/db"
	"gm3u8der/holder"
	"gm3u8der/thm"
	"gm3u8der/ui"
)

const (
	appID = "com.yhyzgn.gm3u8der"
	title = "gm3u8der"
)

func main() {
	db.Init()
	holder.Init()

	der := app.NewWithID(appID)

	// 设置主题
	der.Settings().SetTheme(thm.NewFontTheme())

	// 主窗体
	main := der.NewWindow(title)
	// 设置图标
	ui.Icon(main)

	// 设置主窗体
	ui.Body(main)

	// 设置大小
	main.Resize(fyne.NewSize(1000, 600))
	main.SetFixedSize(true)
	main.SetPadded(false)
	// 显示
	main.ShowAndRun()
}
