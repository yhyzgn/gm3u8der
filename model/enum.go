// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-01 11:33
// version: 1.0.0
// desc   :

package model

import (
	"strings"
)

type ExtType int
type Status int

const (
	MP4 ExtType = iota
	MKV
	AVI
	TS
)

const (
	Downloading Status = iota
	Paused
	Merging
	Finished
	Deleted
)

func MapExtType(ext ExtType) string {
	switch ext {
	case MP4:
		return "MP4"
	case MKV:
		return "MKV"
	case AVI:
		return "AVI"
	default:
		return "TS"
	}
}

func ParseExtType(ext string) ExtType {
	switch strings.ToUpper(ext) {
	case "MP4":
		return MP4
	case "MKV":
		return MKV
	case "AVI":
		return AVI
	default:
		return TS
	}
}
