package common

import (
	"os"
	"time"
)

type Exists interface {
	FileIsExists() bool
}

type Info struct {
	file string
}

func NewExister(file string) *Info {
	return &Info{
		file: file,
	}
}

func (i *Info) FileIsExists() bool {
	fileInfo, err := os.Stat(i.file)
	if err != nil {

		return false
	}
	if fileInfo.ModTime().Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return true
	}
	return false
}
