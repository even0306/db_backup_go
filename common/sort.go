package common

import (
	"db_backup_go/logging"
	"os"
	"sort"
)

// 传入 []os.FileInfo 类型的文件夹，返回文件夹内内容经过排序，由新到旧的同类型文件夹
func SortByTimeFromDirEntry(fs []os.DirEntry) []os.DirEntry {
	sort.SliceStable(fs, func(i, j int) bool {
		fsi, err := fs[i].Info()
		if err != nil {
			logging.Logger.Panic(err)
		}
		fsj, err := fs[j].Info()
		if err != nil {
			logging.Logger.Panic(err)
		}
		flag := false

		if fsi.ModTime().After(fsj.ModTime()) {
			flag = true
		} else if fsi.ModTime().Equal(fsj.ModTime()) {
			if fs[i].Name() < fs[j].Name() {
				flag = true
			}
		}
		return flag
	})
	return fs
}

func SortByTimeFromFileInfo(fs []os.FileInfo) []os.FileInfo {
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
