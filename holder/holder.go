// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 17:24
// version: 1.0.0
// desc   :

package holder

import (
	_ "gm3u8der/env"
	"gm3u8der/model"
	"time"
)

var (
	Settings *model.Settings
)

func Init() {
	// 等待数据库连接
	time.Sleep(1 * time.Second)
	Settings = model.NewSettings().Load()
}
