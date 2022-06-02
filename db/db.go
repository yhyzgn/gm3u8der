// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 22:50
// version: 1.0.0
// desc   :

package db

import (
	"gm3u8der/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"sync"
)

const (
	storeDir   = "./.store"
	dbFilename = storeDir + "/sqlite.db"
)

var (
	once     sync.Once
	Instance *gorm.DB
)

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
