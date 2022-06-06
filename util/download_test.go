// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-06 14:14
// version: 1.0.0
// desc   :

package util

import (
	"path"
	"testing"
)

func TestDownload(t *testing.T) {
	url := "https://github.91chi.fun/https://raw.githubusercontent.com/yhyzgn/gm3u8der/main/ffmpeg.win"
	filename := "ffmpeg.exe"

	Download(path.Join("./", filename), url)
}
