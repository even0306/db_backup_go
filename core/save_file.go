package core

import (
	"io/ioutil"
	"log"
	"mysql_backup_go/config"
)

type save interface {
	SaveFile(f []byte) error
}

type File struct {
}

func (file *File) SaveFile(f []byte, filepath string, filename string) {
	l := config.Logger{}
	logfile, logs := l.SetLogConfig("server.log")
	log.SetOutput(logfile)
	err := ioutil.WriteFile(filepath+filename, f, 0666)
	if err != nil {
		logs.ErrorLogger.Panicf("保存备份文件失败：%v", err)
	}

}
