// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-05-31 10:35
// version: 1.0.0
// desc   :

package util

import (
	"errors"
	"fmt"
	"os"
)

// FileExists 判断文件是否存在
func FileExists(filename string) (exists bool, err error) {
	_, err = os.Stat(filename)
	if nil == err {
		exists = true
		return
	}
	if os.IsNotExist(err) {
		exists = false
		err = errors.New(fmt.Sprintf("no such file '%s'", filename))
		return
	}
	exists = true
	err = errors.New(fmt.Sprintf("the file '%s' is exists, but can not be opened", filename))
	return
}
