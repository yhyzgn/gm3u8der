// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 11:32
// version: 1.0.0
// desc   :

package model

import (
	"fmt"
	"github.com/yhyzgn/goat/file"
	"github.com/yhyzgn/m3u8/model"
	"gm3u8der/dl"
	"gm3u8der/util"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	ffmpeg = "ffmpeg" // ffmpeg 命令
)

type MediaItem struct {
	URL        string
	Name       string
	ExtType    ExtType
	Status     Status
	Progress   float64
	Speed      string
	Downloader *dl.Downloader
}

func (mi *MediaItem) ExtName() string {
	return "." + strings.ToLower(MapExtType(mi.ExtType))
}

func (mi *MediaItem) Download(saveDir string, resolutionSelector func(playList []model.PlayItem, d *dl.Downloader)) {
	tsDir := path.Join(saveDir, "ts_"+util.SHA1(mi.URL))
	mediaFile := mi.Name + mi.ExtName()
	mediaPath := path.Join(saveDir, mediaFile)

	mi.Downloader = dl.New(mi.URL, tsDir)

	go func() {
		for {
			time.Sleep(time.Second)
			mi.Speed = mi.Downloader.Speed()
			mi.Progress = mi.Downloader.Progress()
		}
	}()

	tsNames, tsFile := mi.Downloader.Start(resolutionSelector)

	// 闭包注意传参数，否则引用无效
	go func(names []string, tsF string) {
		for {
			// 下载完再执行
			completed := <-mi.Downloader.Complete()
			if completed {
				// 下载完成，开始合并
				log.Println("TS files download finished, now merging...")
				// 合并切片
				if true {
					// 并转换视频格式
					mergeByFfmpeg(tsDir, mediaPath, mediaFile, tsF)
				} else {
					// 仅合并
					merge(tsDir, mediaPath, mediaFile, names)
				}
			}
		}
	}(tsNames, tsFile)
}

// 合并切片，并转换视频格式
func mergeByFfmpeg(tsDir, mediaPath, mediaFile, tsFile string) {
	// ffmpeg -i "xxx.txt" -acodec copy -vcodec copy -absf aac_adtstoasc out.mp4
	cmdArgs := []string{"-y", "-f", "concat", "-i", tsFile, "-acodec", "copy", "-vcodec", "copy", "-absf", "aac_adtstoasc", mediaPath}

	cmd := exec.Command(ffmpeg, cmdArgs...)

	if err := cmd.Run(); err == nil {
		// 合并完成，删除ts目录
		err := os.RemoveAll(tsDir)
		if nil != err {
			fmt.Println(err)
		} else {
			//fmt.Println(fmt.Sprintf("Media '%s' Merge Finished", colorful(mediaFile)))
		}
	} else {
		fmt.Println(err)
	}
}

// 直接合并成ts文件
func merge(tsDir, mediaPath, mediaFile string, tsNames []string) {
	if nil == tsNames {
		return
	}

	var err error
	for _, name := range tsNames {
		tsFile := path.Join(tsDir, name[5:])
		bs, e := file.Read(tsFile)
		if nil != e {
			err = e
			break
		}
		e = file.Append(mediaPath, bs)
		if nil != e {
			err = e
			break
		}
	}
	if nil == err {
		// 合并完成，删除ts目录
		err := os.RemoveAll(tsDir)
		if nil != err {
			fmt.Println(err)
		} else {
			//fmt.Println(fmt.Sprintf("Media '%s' Merge Finished", colorful(mediaFile)))
		}
	} else {
		fmt.Println(err)
	}
}
