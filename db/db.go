// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 22:50
// version: 1.0.0
// desc   :

package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

const (
	dbFilename = "./assets/db/sqlite.db"
)

var (
	once     sync.Once
	Instance *gorm.DB
)

func Init() {
	once.Do(func() {
		db, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{})
		if nil != err {
			panic(err)
		}
		Instance = db
	})
}
