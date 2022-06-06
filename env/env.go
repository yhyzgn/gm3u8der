// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-02 01:23
// version: 1.0.0
// desc   :

package env

import (
	"fmt"
	"fyne.io/fyne/v2"
	"gm3u8der/util"
	"gm3u8der/wgt"
	"os/exec"
	"path"
	"runtime"
)

const (
	ffmpeg    = "ffmpeg" // ffmpeg 命令
	ffmpegDir = "./"     // ffmpeg 可执行程序所在目录，为达到环境变量优先，这里设置为当前程序运行的目录
	windows   = "https://github.91chi.fun/https://raw.githubusercontent.com/yhyzgn/gm3u8der/main/ffmpeg.win"
	mac       = "https://github.91chi.fun/https://raw.githubusercontent.com/yhyzgn/gm3u8der/main/ffmpeg.mac"
	linux     = "https://github.91chi.fun/https://raw.githubusercontent.com/yhyzgn/gm3u8der/main/ffmpeg.lnx"
)

func Check(win fyne.Window) {
	// 预先检查程序是否存在
	if _, err := exec.LookPath(ffmpeg); nil == err {
		wgt.ShowErrorDialog("错误", "ffmpeg组件缺失，是否下载安装？", win, dispatchFfmpegDownload)
	}
}

// 各系统下载 ffmpeg 组件
func dispatchFfmpegDownload() {
	// 下载 ffmpeg
	gs := runtime.GOOS
	switch gs {
	case "windows":
		ffmpegDownload(windows, ffmpeg+".exe")
	case "darwin":
		ffmpegDownload(mac, ffmpeg)
		if err := exec.Command("chmod", "+X", ffmpeg).Run(); nil != err {
			panic(err)
		}
	case "linux":
		ffmpegDownload(linux, ffmpeg)
		if err := exec.Command("chmod", "+X", ffmpeg).Run(); nil != err {
			panic(err)
		}
	default:
		fmt.Println("Unknown os: ", gs)
	}
}

// 下载 ffmpeg
func ffmpegDownload(url, name string) {
	util.Download(path.Join(ffmpegDir, name), url)
}
