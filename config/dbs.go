package config

import (
	"bufio"
	"log"
	"os"
)

type readDbs interface {
	Read(f string) ([]string, error)
}

type Databases struct {
}

//读取数据库名
func (d *Databases) Read(f string) []string {
	//日志初始化
	l := Logger{}
	logfile, logs := l.SetLogConfig("server.log")
	log.SetOutput(logfile)

	//打开文件
	dbsFile, err := os.OpenFile(f, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		logs.ErrorLogger.Panicf("创建数据库列表文件失败：%v", err)
	}
	defer dbsFile.Close()

	//按行读取文件并添加到切片
	var dbs = []string{}
	line := bufio.NewReader(dbsFile)
	for {
		content, _, err := line.ReadLine()

		if err != nil && content != nil {
			logs.ErrorLogger.Panicf("读取数据库列表行失败：%v", err)
		}
		if content == nil {
			return dbs
		}
		dbs = append(dbs, string(content))
	}
}
