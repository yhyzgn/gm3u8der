// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 11:32
// version: 1.0.0
// desc   :

package model

type MediaItem struct {
	URL      string
	Name     string
	ExtType  ExtType
	SHA1     string
	Status   Status
	Progress float64
	Speed    string
}
