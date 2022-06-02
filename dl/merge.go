// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-02 10:00
// version: 1.0.0
// desc   :

package dl

import (
	"fmt"
	"gm3u8der/util"
	"log"
	"os"
	"os/exec"
	"path"
)

const (
	ffmpeg = "ffmpeg" // ffmpeg 命令
)

type merger struct {
	tsDir         string
	mediaPath     string
	mediaFilename string
	names         []string
	manifest      string
	convert       bool
	finished      bool
}

func (m *merger) apply() {
	// 合并切片
	if m.convert {
		// 并转换视频格式
		m.mergeByFfmpeg()
	} else {
		// 仅合并
		m.merge()
	}
}

// 合并切片，并转换视频格式
func (m *merger) mergeByFfmpeg() {
	// 组装命令参数：ffmpeg -i "xxx.txt" -acodec copy -vcodec copy -absf aac_adtstoasc out.mp4
	cmdArgs := []string{"-y", "-f", "concat", "-i", m.manifest, "-acodec", "copy", "-vcodec", "copy", "-absf", "aac_adtstoasc", m.mediaPath}

	cmd := exec.Command(ffmpeg, cmdArgs...)

	if err := cmd.Run(); err == nil {
		// 合并完成，删除ts目录
		err := os.RemoveAll(m.tsDir)
		if nil != err {
			log.Println(err)
		} else {
			m.finished = true
		}
	} else {
		fmt.Println(err)
	}
}

// 直接合并成ts文件
func (m *merger) merge() {
	if nil == m.names {
		return
	}

	var err error
	for _, name := range m.names {
		tsFile := path.Join(m.tsDir, name[5:])
		bs, e := util.Read(tsFile)
		if nil != e {
			err = e
			break
		}
		e = util.Append(m.mediaPath, bs)
		if nil != e {
			err = e
			break
		}
	}
	if nil == err {
		// 合并完成，删除ts目录
		err := os.RemoveAll(m.tsDir)
		if nil != err {
			fmt.Println(err)
		} else {
			m.finished = true
		}
	} else {
		fmt.Println(err)
	}
}