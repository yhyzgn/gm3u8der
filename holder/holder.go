// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 17:24
// version: 1.0.0
// desc   :

package holder

import "gm3u8der/model"

var (
	Settings *model.Settings
)

func Init() {
	Settings = model.NewSettings().Load()
}
