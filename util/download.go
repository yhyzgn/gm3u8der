// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-02 01:12
// version: 1.0.0
// desc   :

package util

import (
	"bytes"
	"github.com/yhyzgn/m3u8/net"
	"io"
	"io/ioutil"
	"os"
)

func Download(finalPath, url string) {
	tempPath := finalPath + ".tmp"

	// 创建一个临时文件
	target, err := os.Create(tempPath)
	if nil != err {
		return
	}

	// 开始下载
	bys, err := net.Get(url)
	if nil != err {
		return
	}

	realReader := ioutil.NopCloser(bytes.NewReader(bys))

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
}
