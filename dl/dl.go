// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 23:25
// version: 1.0.0
// desc   :

package dl

import (
	"bytes"
	"fmt"
	"github.com/yhyzgn/goat/file"
	"github.com/yhyzgn/m3u8/crypt"
	"github.com/yhyzgn/m3u8/model"
	"github.com/yhyzgn/m3u8/net"
	"github.com/yhyzgn/m3u8/parser"
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
	dir             string
	finished        chan *Resource
	url             string
	speed           string
	speedTicker     *time.Ticker
	lastSize        uint64
	currentSize     uint64
	totalCount      int
	downloadedCount int
	complete        chan bool
}

func New(url, dir string) *Downloader {
	concurrent := runtime.NumCPU()

	return &Downloader{
		wg:          &sync.WaitGroup{},
		pool:        make(chan *Resource, concurrent),
		finished:    make(chan *Resource, concurrent),
		dir:         dir,
		concurrent:  concurrent,
		url:         url,
		speedTicker: time.NewTicker(1 * time.Second),
		complete:    make(chan bool),
	}
}

func (dl *Downloader) Progress() float64 {
	if dl.totalCount > 0 {
		return float64(dl.downloadedCount) / float64(dl.totalCount)
	}
	return 0
}

func (dl *Downloader) Speed() string {
	return dl.speed
}

func (dl *Downloader) Start(resolutionSelector func(playList []model.PlayItem, d *Downloader)) (tsNames []string, tsFile string) {
	m3u, err := parser.FromNetwork(dl.url)
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

	go dl.calcSpeed()

	// 开始下载
	return dl.start(m3u.TsList)
}

func (dl *Downloader) Finished() chan *Resource {
	return dl.finished
}

func (dl *Downloader) Complete() chan bool {
	return dl.complete
}

func (dl *Downloader) append(resources ...*Resource) *Downloader {
	dl.resources = append(dl.resources, resources...)
	return dl
}

func (dl *Downloader) start(tsList []model.TS) (tsNames []string, tsFile string) {
	if err := os.MkdirAll(dl.dir, os.ModePerm); nil != err {
		panic(err)
	}

	dl.totalCount = len(tsList)

	keyMap := make(map[string][]byte)
	tsNames = make([]string, 0)

	for i, item := range tsList {
		name := fmt.Sprintf("slice_%06d.ts", i+1)
		tsNames = append(tsNames, "file "+name)

		if nil != item.Key && item.Key.URI != "" && nil == keyMap[string(item.Key.Method)+"-"+item.Key.URI] {
			keyMap[string(item.Key.Method)+"-"+item.Key.URI], _ = net.Get(item.Key.URI)
		}

		dl.append(&Resource{
			index:    i,
			url:      item.URL,
			filename: name,
			override: false,
		})
	}

	tsFile = path.Join(dl.dir, "slice.lst")
	_ = file.WriteString(tsFile, strings.Join(tsNames, "\n"))

	// 更新进度条
	go func() {
		for {
			<-dl.Finished()
			dl.downloadedCount++
			dl.complete <- dl.downloadedCount == dl.totalCount
		}
	}()

	go dl.runWithReader(func(resourceIndex int, reader io.ReadCloser) io.Reader {
		key := tsList[resourceIndex].Key
		if nil == key {
			return reader
		}
		data, _ := ioutil.ReadAll(reader)
		data, _ = crypt.AES128Decrypt(data, keyMap[string(key.Method)+"-"+key.URI], []byte(key.IV))
		return bytes.NewReader(data)
	})
	return
}

func (dl *Downloader) runWithReader(reader func(resourceIndex int, reader io.ReadCloser) io.Reader) {
	for i, resource := range dl.resources {
		dl.wg.Add(1)
		resource.index = i
		go dl.download(resource, reader)
	}
	dl.wg.Wait()

	// 速度计算定时器停止
	dl.speedTicker.Stop()
}

func (dl *Downloader) run() {
	dl.runWithReader(func(resourceIndex int, reader io.ReadCloser) io.Reader {
		return reader
	})
}

func (dl *Downloader) download(resource *Resource, reader func(resourceIndex int, reader io.ReadCloser) io.Reader) {
	defer dl.wg.Done()
	dl.pool <- resource
	finalPath := path.Join(dl.dir, resource.filename)

	// 如果不覆盖下载，文件存在时则无需下载
	if file.Exists(finalPath) && !resource.override {
		// 也表示完成一个任务
		dl.finished <- <-dl.pool
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
	dl.finished <- <-dl.pool
}

func (dl *Downloader) calcSpeed() {
	go func() {
		for _ = range dl.speedTicker.C {
			delta := dl.currentSize - dl.lastSize
			dl.speed = util.FormatFileSize(delta) + "/s"
			dl.lastSize = dl.currentSize
		}
	}()
}
