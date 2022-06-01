// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 12:39
// version: 1.0.0
// desc   :

package util

func If[T any](condition bool, positive, negative T) T {
	if condition {
		return positive
	}
	return negative
}
