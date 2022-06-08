// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 11:32
// version: 1.0.0
// desc   :

package model

import (
	"github.com/yhyzgn/m3u8/model"
	"gm3u8der/dl"
	"strings"
)

type MediaItem struct {
	URL        string
	Name       string
	ExtType    ExtType
	Status     Status
	Downloader *dl.Downloader
}

func (mi *MediaItem) ExtName() string {
	return "." + strings.ToLower(MapExtType(mi.ExtType))
}

func (mi *MediaItem) Download(saveDir string, resolutionSelector func(playList []model.PlayItem, d *dl.Downloader)) {
	mi.Downloader = dl.New(mi.URL, saveDir, mi.Name, mi.ExtName())
	mi.Downloader.Start(resolutionSelector)
}

func (mi *MediaItem) Speed() string {
	if nil == mi.Downloader {
		return "正在准备..."
	}
	return mi.Downloader.Speed()
}

func (mi *MediaItem) Progress() float64 {
	if nil == mi.Downloader {
		return 0
	}
	progress := mi.Downloader.Progress()
	if progress == 1 {
		mi.Status = Finished
	} else if progress >= 0.96 {
		mi.Status = Merging
	} else {
		mi.Status = Downloading
	}
	return progress
}
