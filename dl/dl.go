// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 23:25
// version: 1.0.0
// desc   :

package dl

import (
	"bytes"
	"fmt"
	"github.com/yhyzgn/m3u8/crypt"
	"github.com/yhyzgn/m3u8/model"
	"github.com/yhyzgn/m3u8/net"
	"github.com/yhyzgn/m3u8/parser"
	"gm3u8der/component"
	"gm3u8der/util"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Downloader struct {
	wg              *sync.WaitGroup
	pool            chan *Resource
	concurrent      int
	resources       []*Resource
	tsDir           string
	resourceFinish  chan *Resource
	m3u8URL         string
	lastSize        uint64
	currentSize     uint64
	totalCount      int
	downloadedCount int
	started         bool
	allFinished     bool
	progress        float64
	speed           string
	merger          *merger
	calcTicker      *component.TimeTicker
}

func New(m3u8URL, saveDir, name, extName string) *Downloader {
	tsDir := path.Join(saveDir, "ts_"+name)
	mediaFile := name + extName
	mediaPath := path.Join(saveDir, mediaFile)

	concurrent := runtime.NumCPU()

	return &Downloader{
		wg:             &sync.WaitGroup{},
		pool:           make(chan *Resource, concurrent),
		resourceFinish: make(chan *Resource, concurrent),
		tsDir:          tsDir,
		concurrent:     concurrent,
		m3u8URL:        m3u8URL,
		merger: &merger{
			tsDir:         tsDir,
			mediaFilePath: mediaPath,
			mediaFilename: mediaFile,
			convert:       ".ts" != extName,
		},
	}
}

func (dl *Downloader) Speed() string {
	return dl.speed
}

func (dl *Downloader) Progress() float64 {
	return dl.progress
}

func (dl *Downloader) Start(resolutionSelector func(playList []model.PlayItem, d *Downloader), overrideTip func(mediaFilename string, override chan bool), onStarted func()) {
	m3u, err := parser.FromNetwork(dl.m3u8URL)
	if nil != err {
		log.Println(err)
		return
	}

	if len(m3u.PlayList) > 0 {
		// 选择分辨率方案
		resolutionSelector(m3u.PlayList, dl)
		return
	}

	if len(m3u.TsList) == 0 {
		log.Println("未获取到任何 ts 片段信息")
		return
	}

	// 检查媒体文件是否已存在
	if exists, _ := util.FileExists(dl.merger.mediaFilePath); exists {
		override := make(chan bool)
		overrideTip(dl.merger.mediaFilename, override)
		// 接收选择结果
		go func(innerM3u *model.M3U8, ovr chan bool, started func()) {
			for {
				// 不覆盖，就不下载啦~
				if <-ovr {
					// 开始覆盖下载
					dl.start(innerM3u.TsList, started)
					break
				}
			}
		}(m3u, override, onStarted)
	} else {
		// 开始下载
		dl.start(m3u.TsList, onStarted)
	}
}

func (dl *Downloader) append(resources ...*Resource) *Downloader {
	dl.resources = append(dl.resources, resources...)
	return dl
}

func (dl *Downloader) start(tsList []model.TS, onStarted func()) {
	if onStarted != nil {
		onStarted()
	}

	if err := os.MkdirAll(dl.tsDir, os.ModePerm); nil != err {
		panic(err)
	}

	dl.totalCount = len(tsList)
	dl.started = true

	keyMap := make(map[string][]byte)
	tsNames := make([]string, 0)

	for i, item := range tsList {
		name := fmt.Sprintf("slice_%06d.ts", i+1)
		tsNames = append(tsNames, "file "+name)

		if nil != item.Key && item.Key.URI != "" {
			// TODO 本次的 key 不带参数。。
			item.Key.URI = strings.Split(item.Key.URI, "?")[0]
			if nil == keyMap[string(item.Key.Method)+"-"+item.Key.URI] {
				keyMap[string(item.Key.Method)+"-"+item.Key.URI], _ = net.Get(item.Key.URI)
			}
		}

		dl.append(&Resource{
			index:    i,
			url:      item.URL,
			filename: name,
			override: false,
		})
	}

	tsFile := path.Join(dl.tsDir, "slice.lst")
	_ = util.WriteString(tsFile, strings.Join(tsNames, "\n"))

	// 更新进度
	go func(names []string, manifestFile string) {
		for {
			<-dl.resourceFinish
			dl.downloadedCount++
			if dl.downloadedCount == dl.totalCount {
				dl.allFinished = true

				dl.merger.names = names
				dl.merger.manifest = manifestFile

				// 开始合并视频片段
				dl.merger.apply(func() {
					// 下载完成后，停止计时器
					if nil != dl.calcTicker {
						dl.calcTicker.Stop()
						dl.calcTicker = nil
					}
				})
			}
		}
	}(tsNames, tsFile)

	go dl.runWithReader(func(resourceIndex int, reader io.ReadCloser) io.Reader {
		key := tsList[resourceIndex].Key
		if nil == key {
			return reader
		}
		data, _ := ioutil.ReadAll(reader)
		data, _ = crypt.AES128Decrypt(data, keyMap[string(key.Method)+"-"+key.URI], []byte(key.IV))
		return bytes.NewReader(data)
	})

	// 计算速度和进度
	go func() {
		dl.calcTicker = component.StartTicker(1*time.Second, func() {
			dl.calcSpeed()
			dl.calcProgress()
		})
	}()
}

func (dl *Downloader) runWithReader(reader func(resourceIndex int, reader io.ReadCloser) io.Reader) {
	for i, resource := range dl.resources {
		dl.wg.Add(1)
		resource.index = i
		go dl.download(resource, reader)
	}
	dl.wg.Wait()
}

func (dl *Downloader) run() {
	dl.runWithReader(func(resourceIndex int, reader io.ReadCloser) io.Reader {
		return reader
	})
}

func (dl *Downloader) download(resource *Resource, reader func(resourceIndex int, reader io.ReadCloser) io.Reader) {
	defer dl.wg.Done()
	dl.pool <- resource
	finalPath := path.Join(dl.tsDir, resource.filename)

	// 如果不覆盖下载，文件存在时则无需下载
	if exists, err := util.FileExists(finalPath); nil == err && exists && !resource.override {
		// 也表示完成一个任务
		dl.resourceFinish <- <-dl.pool
		return
	}

	tempPath := finalPath + ".tmp"

	// 创建一个临时文件
	target, err := os.Create(tempPath)
	if nil != err {
		return
	}

	// 开始下载
	bys, err := net.Get(resource.url)
	if nil != err {
		return
	}

	// 记录最新大小
	dl.currentSize += uint64(len(bys))

	realReader := reader(resource.index, ioutil.NopCloser(bytes.NewReader(bys)))

	// 将下载的文件流写到临时文件
	_, err = io.Copy(target, realReader)
	if nil != err {
		_ = target.Close()
		return
	}

	_ = target.Close()
	err = os.Rename(tempPath, finalPath)
	if nil != err {
		return
	}

	// 完成一个任务
	dl.resourceFinish <- <-dl.pool
}

func (dl *Downloader) calcSpeed() {
	if !dl.started {
		dl.speed = "正在准备..."
		return
	}
	if dl.merger.finished {
		dl.speed = "下载完成"
		return
	}
	if dl.allFinished {
		dl.speed = "正在合并..."
		return
	}
	// 下载中
	deltaSize := dl.currentSize - dl.lastSize
	dl.lastSize = dl.currentSize
	dl.speed = util.FormatFileSize(deltaSize) + "/s"
}

func (dl *Downloader) calcProgress() {
	if dl.totalCount > 0 {
		// 下载过程占 96%，合并过程占 4%
		downloadProgress := float64(dl.downloadedCount) / float64(dl.totalCount)
		if !dl.allFinished {
			// 下载未结束
			dl.progress = downloadProgress * 0.96
			return
		}
		if dl.allFinished && dl.merger.finished {
			// 已经合并完成
			dl.progress = 1.0
			return
		}
		dl.progress = downloadProgress
		return
	}
	dl.progress = 0
}
