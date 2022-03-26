package controller

import (
	"mysql_backup_go/config"
	"mysql_backup_go/core"
)

type Controller interface {
	Controller() error
}

type FileInfo struct {
}

func NewController() *FileInfo {
	return &FileInfo{}
}

//备份主程序
func (bk FileInfo) Controller() error {
	//获取配置文件
	conf := config.NewConfig("config.json")
	confData, err := conf.Read()
	if err != nil {
		return err
	}
	//获取要使用的数据库列表
	dbs := config.NewDBList("dbs.txt")
	dbsData, err := dbs.Read()
	if err != nil {
		return err
	}

	cp := core.NewCompartor(confData, &dbsData)
	preDBS, err := cp.Comparison()
	if err != nil {
		return err
	}

	//开始循环备份每个数据库
	for _, v := range *preDBS {
		bk := core.NewBackuper(confData, &v)
		go bk.Run()
	}

	//按天保留最新7份备份，删除之前的备份

	return nil
}
