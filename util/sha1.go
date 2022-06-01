// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 12:42
// version: 1.0.0
// desc   :

package util

import (
	"crypto/sha1"
	"encoding/hex"
)

func SHA1(data string) string {
	o := sha1.New()
	o.Write([]byte(data))
	return hex.EncodeToString(o.Sum(nil))
}
