// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 17:31
// version: 1.0.0
// desc   :

package util

import (
	"os"
	"path/filepath"
)

func SystemDownloadDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, "Downloads")
}
