package common

import (
	"os"
	"sort"
)

type Sort interface {
	SortByTime() []os.FileInfo
}

type Files struct {
	files []os.FileInfo
}

func NewOrder(files []os.FileInfo) *Files {
	return &Files{
		files: files,
	}
}

func (fs *Files) SortByTime() []os.FileInfo {
	sort.SliceStable(fs.files, func(i, j int) bool {
		flag := false
		if fs.files[i].ModTime().After(fs.files[j].ModTime()) {
			flag = true
		} else if fs.files[i].ModTime().Equal(fs.files[j].ModTime()) {
			if fs.files[i].Name() < fs.files[j].Name() {
				flag = true
			}
		}
		return flag
	})
	return fs.files
}
