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
	"gm3u8der/component"
	"gm3u8der/dl"
	"gm3u8der/env"
	"gm3u8der/holder"
	"gm3u8der/model"
	"gm3u8der/util"
	"gm3u8der/wgt"
	"log"
	"strconv"
	"time"
)

var (
	divLine = canvas.NewLine(theme.ShadowColor())

	listDownloading   *widget.List
	bdDownLoadingList = make([]*model.MediaItem, 0)

	dlgURLShowing = false // 弹窗是否正在显示
)

func init() {
	divLine.StrokeWidth = 1
}

// Body ...
func Body(win fyne.Window) {
	// 监视剪贴板
	Clipboard(win, func(m3u8URL string, onShowed func()) {
		if !dlgURLShowing {
			dlgURLShowing = true
			showURLDialog(win, m3u8URL, func() {
				dlgURLShowing = false
			})

			onShowed()
		}
	})

	// 创建工具栏
	toolBar := widget.NewToolbar(
		widget.NewToolbarAction(resourceIcaddPng, func() {
			if !dlgURLShowing {
				dlgURLShowing = true
				showURLDialog(win, "", func() {
					dlgURLShowing = false
				})
			}
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(resourceIcsettingsPng, func() {
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

			wdCheckClipboard := widget.NewCheck("是否开启", func(b bool) {})
			wdCheckClipboard.Checked = holder.Settings.Clipboard
			ckClipboard := widget.NewFormItem("监听剪切板", wdCheckClipboard)

			rgTraySticky := widget.NewRadioGroup([]string{"固定到任务栏", "直接退出应用"}, func(s string) {})
			rgTraySticky.SetSelected(func() string {
				if holder.Settings.TraySticky {
					return "固定到任务栏"
				}
				return "直接退出应用"
			}())
			rgTraySticky.Horizontal = true
			fiTraySticky := widget.NewFormItem("关闭窗口时", rgTraySticky)

			items := []*widget.FormItem{itemSelectDir, fiExt, taskCount, ckClipboard, fiTraySticky}
			var dlg dialog.Dialog
			dlg = dialog.NewForm("设置", "保存", "取消", items, func(b bool) {
				if b {
					// 保存
					taskCount, err := strconv.Atoi(wdTaskCount.Text)
					if nil != err {
						wgt.ShowErrorDialog("任务数必须是数字", "是否返回输入？", win, func() {
							taskCount = holder.Settings.TaskCount
							dlg.Show()
						}, nil)
						return
					}
					holder.Settings.SaveDir = wdSaveDir.Text
					holder.Settings.ExtType = model.ParseExtType(rg.Selected)
					holder.Settings.TaskCount = taskCount
					holder.Settings.Clipboard = wdCheckClipboard.Checked
					holder.Settings.TraySticky = rgTraySticky.Selected == "固定到任务栏"

					// 保存
					holder.Settings.Store()
				}
			}, win)
			dlg.Resize(fyne.Size{Width: 600})
			dlg.Show()
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {}),
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
		info.Objects[2].(*widget.Label).SetText(item.Speed())

		progress.SetValue(item.Progress())
	})

	banner := container.NewVBox(toolBar, divLine)
	content := container.New(
		layout.NewBorderLayout(banner, nil, nil, nil),
		banner,
		listDownloading,
	)

	win.SetContent(content)

	// UI设置后，检查环境
	env.Check(win)

	// 定时刷新列表
	component.StartTicker(time.Second*1, func() {
		if nil != listDownloading && nil != bdDownLoadingList && len(bdDownLoadingList) > 0 {
			listDownloading.Refresh()
		}
	})
}

func showURLDialog(win fyne.Window, providedURL string, onClosed func()) {
	rg := widget.NewRadioGroup([]string{"MP4", "MKV", "AVI", "TS"}, func(s string) {})
	rg.SetSelected(model.MapExtType(holder.Settings.ExtType))
	rg.Horizontal = true

	// URL
	wdURL := widget.NewEntry()
	wdURL.SetText(providedURL)

	// 媒体名称
	wdName := widget.NewEntry()
	if "" != providedURL {
		wdName.SetText(util.SHA1(providedURL))
	}

	fiM3u8URL := &widget.FormItem{
		Text:   "M3U8链接",
		Widget: wdURL,
	}

	fiMediaName := &widget.FormItem{
		Text:   "媒体名称",
		Widget: wdName,
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
				wgt.ShowErrorDialog("m3u8地址错误", "是否重新输入？", win, func() {
					dlg.Show()
				}, nil)
				return
			}
			valMediaName := fiMediaName.Widget.(*widget.Entry).Text
			if valMediaName == "" {
				wgt.ShowErrorDialog("媒体文件名称为空", "是否返回输入？", win, func() {
					dlg.Show()
				}, nil)
				return
			}

			valExt := fiExt.Widget.(*widget.RadioGroup).Selected

			item := &model.MediaItem{
				URL:     valM3u8URL,
				Name:    valMediaName,
				ExtType: model.ParseExtType(valExt),
				Status:  model.Downloading,
			}

			// 试试看
			var selector func(playList []m3u8Model.PlayItem, d *dl.Downloader)
			selector = func(playList []m3u8Model.PlayItem, d *dl.Downloader) {
				// 选择分辨率方案
				rgItems := make([]string, len(playList))
				itemIndexMap := make(map[string]int)

				var label string
				for i, pl := range playList {
					if "" != pl.Resolution {
						label = "分辨率"
						rgItems[i] = pl.Resolution
					} else {
						label = "带宽"
						rgItems[i] = pl.BandWidth
					}
					itemIndexMap[rgItems[i]] = i
				}
				rgResolution := widget.NewRadioGroup(rgItems, func(s string) {})
				rgResolution.SetSelected(rgItems[len(rgItems)-1]) // 默认选中最后一个

				fiResolution := &widget.FormItem{
					Text:   label,
					Widget: rgResolution,
				}

				dialog.ShowForm("选择视频分辨率", "确定", "取消", []*widget.FormItem{fiResolution}, func(bOk bool) {
					if bOk {
						// 确定
						piSelected := playList[itemIndexMap[rgResolution.Selected]]
						item.URL = piSelected.URL
						item.Download(holder.Settings.SaveDir, selector, func(mediaFilename string, override chan bool) {
							wgt.ShowErrorDialog("警告", "文件【"+mediaFilename+"】已存在，是否覆盖下载？", win, func() {
								override <- true
							}, func() {
								override <- false
							})
						}, func() {
							// 开始后将媒体添加到列表
							bdDownLoadingList = append([]*model.MediaItem{item}, bdDownLoadingList...)
						})
					}
					// 取消就不管啦
				}, win)
			}
			item.Download(holder.Settings.SaveDir, selector, func(mediaFilename string, override chan bool) {
				wgt.ShowErrorDialog("警告", "文件【"+mediaFilename+"】已存在，是否覆盖下载？", win, func() {
					override <- true
				}, func() {
					override <- false
				})
			}, func() {
				// 开始后将媒体添加到列表
				bdDownLoadingList = append([]*model.MediaItem{item}, bdDownLoadingList...)
			})
		} else {
			// 取消
		}
	}, win)
	dlg.Resize(fyne.Size{
		Width: 700,
	})
	if nil != onClosed {
		dlg.SetOnClosed(onClosed)
	}
	dlg.Show()
}
