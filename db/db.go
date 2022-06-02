// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 22:50
// version: 1.0.0
// desc   :

package db

import (
	"gm3u8der/cst"
	"gm3u8der/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"sync"
)

const (
	AppID = "com.yhyzgn.gm3u8der"
)

var (
	once       sync.Once
	storeDir   string
	dbFilename string
	Instance   *gorm.DB
)

func init() {
	configDir, err := os.UserConfigDir()
	if nil != err {
		panic(err)
	}
	storeDir = path.Join(configDir, cst.AppID)
	dbFilename = path.Join(storeDir, "sqlite.db")
}

func Init() {
	once.Do(func() {
		if exists, err := util.FileExists(storeDir); nil != err || !exists {
			err = os.MkdirAll(storeDir, os.ModePerm)
			if nil != err {
				panic(err)
			}
		}
		if exists, err := util.FileExists(dbFilename); nil != err || !exists {
			dbFile, err := os.Create(dbFilename)
			if nil != err {
				panic(err)
			}
			_ = dbFile.Close()
		}

		db, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{})
		if nil != err {
			panic(err)
		}
		Instance = db
	})
}
