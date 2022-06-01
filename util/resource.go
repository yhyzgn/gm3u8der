// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 9:53
// version: 1.0.0
// desc   :

package util

import (
	"fmt"
	"fyne.io/fyne/v2"
)

// LoadResourceFromFile 从文件加载资源
func LoadResourceFromFile(filename string) fyne.Resource {
	exists, err := FileExists(filename)
	if !exists || nil != err {
		panic(err)
	}
	res, err := fyne.LoadResourceFromPath(filename)
	if nil != err {
		panic(fmt.Sprintf("The file '%s' loading failed.", filename))
	}
	return res
}
