// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-02 15:58
// version: 1.0.0
// desc   :

package ui

import "testing"

func TestClipboard(t *testing.T) {
	t.Log(m3u8URLRegexp.MatchString("http://devimages.apple.com/iphone/samples/bipbop/gear1/prog_index.m3u8"))
	t.Log(m3u8URLRegexp.MatchString("https://devimages.apple.com/iphone/samples/bipbop/gear1/prog_index.m3u8"))
	t.Log(m3u8URLRegexp.MatchString("https://devimages.apple.com/iphone/samples/bipbop/gear1/prog_index.m3u8?token=1234&timestamp=12345"))
	t.Log(m3u8URLRegexp.MatchString("https://devimages.apple.com/iphone/samples/bipbop/gear1/prog_index.asc"))
}
