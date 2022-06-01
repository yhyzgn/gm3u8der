// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 12:46
// version: 1.0.0
// desc   :

package util

import "testing"

func TestSHA1(t *testing.T) {
	txt := "Hello 你好"
	res := SHA1(txt)
	t.Log(res)
}
