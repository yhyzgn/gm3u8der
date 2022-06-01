// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 11:24
// version: 1.0.0
// desc   :

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
	m3u8Model "github.com/yhyzgn/m3u8/model"
	"gm3u8der/dl"
	"gm3u8der/holder"
	"gm3u8der/model"
	"gm3u8der/util"
	"log"
	"strconv"
	"time"
)

const (
	icAdd      = "./assets/img/ic_add.png"
	icDelete   = "./assets/img/ic_delete.png"
	icSettings = "./assets/img/ic_settings.png"
)

var (
	divLine = canvas.NewLine(theme.ShadowColor())
)

func init() {
	divLine.StrokeWidth = 1
}

func Body(win fyne.Window) {
	var listDownloading *widget.List
	bdDownLoadingList := make([]*model.MediaItem, 0)

	toolBar := widget.NewToolbar(
		widget.NewToolbarAction(util.LoadResourceFromFile(icAdd), func() {
			rg := widget.NewRadioGroup([]string{"MP4", "MKV", "AVI", "TS"}, func(s string) {})
			rg.SetSelected(model.MapExtType(holder.Settings.ExtType))
			rg.Horizontal = true

			fiM3u8URL := &widget.FormItem{
				Text:   "M3U8链接",
				Widget: widget.NewEntry(),
			}

			fiMediaName := &widget.FormItem{
				Text:   "媒体名称",
				Widget: widget.NewEntry(),
			}
			fiExt := &widget.FormItem{
				Text:   "输出格式",
				Widget: rg,
			}
			items := []*widget.FormItem{fiM3u8URL, fiMediaName, fiExt}
			var dlg dialog.Dialog
			dlg = dialog.NewForm("新建下载任务", "下载", "取消", items, func(b bool) {
				if b {
					// 确定
					valM3u8URL := fiM3u8URL.Widget.(*widget.Entry).Text
					if valM3u8URL == "" {
						dialog.NewConfirm("m3u8地址错误", "是否重新输入？", func(ok bool) {
							if ok {
								dlg.Show()
							}
						}, win).Show()
						return
					}
					valMediaName := fiMediaName.Widget.(*widget.Entry).Text
					if valMediaName == "" {
						dialog.NewConfirm("媒体文件名称为空", "是否返回输入？", func(ok bool) {
							if ok {
								dlg.Show()
							}
						}, win).Show()
						return
					}

					valExt := fiExt.Widget.(*widget.RadioGroup).Selected

					item := &model.MediaItem{
						URL:      valM3u8URL,
						Name:     valMediaName,
						ExtType:  model.ParseExtType(valExt),
						Status:   model.Downloading,
						Progress: 0.0,
					}

					// 试试看
					var fuck func(playList []m3u8Model.PlayItem, d *dl.Downloader)
					fuck = func(playList []m3u8Model.PlayItem, d *dl.Downloader) {
						item.URL = playList[0].URL
						item.Download(holder.Settings.SaveDir, fuck)
					}
					item.Download(holder.Settings.SaveDir, fuck)

					bdDownLoadingList = append([]*model.MediaItem{item}, bdDownLoadingList...)
				} else {
					// 取消
				}
			}, win)
			dlg.Resize(fyne.Size{
				Width: 700,
			})
			dlg.Show()
		}),
		widget.NewToolbarAction(util.LoadResourceFromFile(icDelete), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(util.LoadResourceFromFile(icSettings), func() {
			wdSaveDir := widget.NewEntry()
			wdSaveDir.SetText(holder.Settings.SaveDir)
			wdBtnOpenDir := widget.NewButton("选择", func() {
				dir, err := zenity.SelectFile(zenity.Title("选择下载目录"), zenity.Directory())
				if err != nil {
					log.Println(err)
					return
				}
				wdSaveDir.SetText(dir)
			})
			itemSelectDir := widget.NewFormItem("媒体下载目录", container.NewBorder(nil, nil, nil, wdBtnOpenDir, wdBtnOpenDir, wdSaveDir))

			rg := widget.NewRadioGroup([]string{"MP4", "MKV", "AVI", "TS"}, func(s string) {})
			rg.SetSelected(model.MapExtType(holder.Settings.ExtType))
			rg.Horizontal = true
			fiExt := widget.NewFormItem("默认输出格式", rg)

			wdTaskCount := widget.NewEntry()
			wdTaskCount.SetText(strconv.Itoa(holder.Settings.TaskCount))
			taskCount := widget.NewFormItem("同时下载任务数", wdTaskCount)

			items := []*widget.FormItem{itemSelectDir, fiExt, taskCount}
			dlg := dialog.NewForm("设置", "保存", "取消", items, func(b bool) {
				if b {
					// 保存
					taskCount, err := strconv.Atoi(wdTaskCount.Text)
					if nil != err {
						log.Println("任务数必须是数字")
						taskCount = holder.Settings.TaskCount
					}
					holder.Settings.SaveDir = wdSaveDir.Text
					holder.Settings.ExtType = model.ParseExtType(rg.Selected)
					holder.Settings.TaskCount = taskCount

					// 保存
					holder.Settings.Store()
				}
			}, win)
			dlg.Resize(fyne.Size{Width: 600})
			dlg.Show()
		}),
	)

	listDownloading = widget.NewList(func() int {
		return len(bdDownLoadingList)
	}, func() fyne.CanvasObject {
		pb := widget.NewProgressBar()
		pb.Resize(fyne.Size{Height: 4})

		return container.NewVBox(
			container.NewHBox(
				widget.NewLabel("名称"),
				layout.NewSpacer(),
				widget.NewLabel("网速"),
			),
			pb,
		)
	}, func(id widget.ListItemID, co fyne.CanvasObject) {
		item := bdDownLoadingList[id]

		info := co.(*fyne.Container).Objects[0].(*fyne.Container)
		progress := co.(*fyne.Container).Objects[1].(*widget.ProgressBar)

		name := []rune(item.Name)
		if len(name) > 50 {
			end := name[len(name)-6:]
			name = append(name[:40], []rune("...")...)
			name = append(name, end...)
		}
		info.Objects[0].(*widget.Label).SetText(string(name) + item.ExtName())
		info.Objects[2].(*widget.Label).SetText(item.Speed)

		progress.SetValue(item.Progress)
	})

	banner := container.NewVBox(toolBar, divLine)
	content := container.New(
		layout.NewBorderLayout(banner, nil, nil, nil),
		banner,
		listDownloading,
	)

	go func() {
		for {
			time.Sleep(time.Second)
			if nil != listDownloading {
				listDownloading.Refresh()
			}
		}
	}()

	win.SetContent(content)
}
