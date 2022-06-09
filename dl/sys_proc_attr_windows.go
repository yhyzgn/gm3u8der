// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-09 9:42
// version: 1.0.0
// desc   :

//go:build windows

package dl

import (
	"os/exec"
	"syscall"
)

func applySysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
