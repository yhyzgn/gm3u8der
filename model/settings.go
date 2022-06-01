// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 17:17
// version: 1.0.0
// desc   :

package model

import (
	"encoding/json"
	"gm3u8der/util"
	"io/ioutil"
	"os"
)

const (
	filename = "./settings.json"
)

type Settings struct {
	SaveDir   string  `json:"saveDir"`
	ExtType   ExtType `json:"extType"`
	TaskCount int     `json:"taskCount"`
}

func NewSettings() *Settings {
	return new(Settings)
}

func (s *Settings) Load() *Settings {
	bys, err := ioutil.ReadFile(filename)
	if nil != err || nil == bys {
		s.storeDefault()
	}

	var temp Settings
	err = json.Unmarshal(bys, &temp)
	if nil != err {
		s.storeDefault()
	} else {
		s.SaveDir = temp.SaveDir
		s.ExtType = temp.ExtType
		s.TaskCount = temp.TaskCount
	}
	return s
}

func (s *Settings) Store() {
	s.store()
}

func (s *Settings) storeDefault() {
	s.SaveDir = util.SystemDownloadDir()
	s.ExtType = MP4
	s.TaskCount = 5

	s.store()
}

func (s *Settings) store() {
	bys, err := json.MarshalIndent(s, "", "\t")
	if nil != err {
		panic(err)
	}

	if err = ioutil.WriteFile(filename, bys, os.ModePerm); nil != err {
		panic(err)
	}
}
