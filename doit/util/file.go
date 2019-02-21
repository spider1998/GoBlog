package util

import (
	"os"

	"github.com/pkg/errors"
)

//创建文件
func MakeDirectory(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, os.ModePerm)
		}
		return err
	}
	if !fi.IsDir() {
		return errors.New("specified path is not a directory")
	}
	return nil
}

//判断文件存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
