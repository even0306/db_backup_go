package common

import (
	"os"
	"sort"
)

//传入 []os.FileInfo 类型的文件夹，返回文件夹内内容经过排序，由新到旧的同类型文件夹
func SortByTime(fs []os.FileInfo) []os.FileInfo {
	sort.SliceStable(fs, func(i, j int) bool {
		flag := false
		if fs[i].ModTime().After(fs[j].ModTime()) {
			flag = true
		} else if fs[i].ModTime().Equal(fs[j].ModTime()) {
			if fs[i].Name() < fs[j].Name() {
				flag = true
			}
		}
		return flag
	})
	return fs
}
