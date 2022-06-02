// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-05-31 10:35
// version: 1.0.0
// desc   :

package util

import (
	"errors"
	"fmt"
	"io/ioutil"
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

// WriteString 写字符串到文件
func WriteString(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	return writeString(file, data)
}

// Append 写文件追加
func Append(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	return writeBytes(file, data)
}

func writeBytes(file *os.File, data []byte) error {
	_, err := file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func writeString(file *os.File, data string) error {
	_, err := file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

// Read 读文件
func Read(filename string) (data []byte, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(file)
	_ = file.Close()
	return
}
