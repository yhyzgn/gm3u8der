// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 17:17
// version: 1.0.0
// desc   :

package model

import (
	"gm3u8der/db"
	"gm3u8der/util"
)

// Settings ...
type Settings struct {
	ID         int     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:int(11);"`
	SaveDir    string  `json:"saveDir" gorm:"column:save_dir;type:varchar(255);"`
	ExtType    ExtType `json:"extType" gorm:"column:ext_type;type:int(11);"`
	TaskCount  int     `json:"taskCount" gorm:"column:task_count;type:int(11);"`
	Clipboard  bool    `json:"clipboard" gorm:"column:clipboard;type:tinyint(1);"`
	TraySticky bool    `json:"TraySticky" gorm:"column:tray_sticky;type:tinyint(1);"`
}

// NewSettings ...
func NewSettings() *Settings {
	_ = db.Instance.AutoMigrate(&Settings{})
	return new(Settings)
}

// Load ...
func (s *Settings) Load() *Settings {
	var temp Settings
	err := db.Instance.First(&temp).Error
	if nil != err {
		s.storeDefault()
		return s
	}

	s.ID = temp.ID
	s.SaveDir = temp.SaveDir
	s.ExtType = temp.ExtType
	s.TaskCount = temp.TaskCount
	s.Clipboard = temp.Clipboard
	s.TraySticky = temp.TraySticky

	return s
}

// Store ...
func (s *Settings) Store() {
	s.store()
}

func (s *Settings) storeDefault() {
	s.ID = 1
	s.SaveDir = util.SystemDownloadDir()
	s.ExtType = MP4
	s.TaskCount = 5
	s.Clipboard = true
	s.TraySticky = true

	s.store()
}

func (s *Settings) store() {
	err := db.Instance.Save(s).Error
	if nil != err {
		panic(err)
	}
}
