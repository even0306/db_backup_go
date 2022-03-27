package controller

import (
	"mysql_backup_go/common"
	"mysql_backup_go/modules"
	"sync"
)

type Controller interface {
	Controller() error
}

type fileInfo struct {
	confFile string
	dbsFile  string
	fileName string
}

func NewController(conf string, dbs string) *fileInfo {
	return &fileInfo{
		confFile: conf,
		dbsFile:  dbs,
		fileName: "",
	}
}

//备份主程序
func (fi fileInfo) Controller() error {
	//获取配置文件
	conf := common.NewConfig(fi.confFile)
	confData, err := conf.Read()
	if err != nil {
		return err
	}
	//获取要使用的数据库列表
	dbs := common.NewDBList(fi.dbsFile)
	dbsData, err := dbs.Read()
	if err != nil {
		return err
	}

	cp := modules.NewCompartor(confData, dbsData)
	preDBS, err := cp.Comparison()
	if err != nil {
		return err
	}

	//开始循环备份每个数据库
	var wg sync.WaitGroup
	for _, v := range *preDBS {
		bk := modules.NewBackuper(confData, &v)
		wg.Add(1)
		go func() error {
			fi.fileName, err = bk.Run()
			if err != nil {
				return err
			}

			return nil
		}()
	}
	wg.Wait()

	//按天保留最新7份备份，删除之前的备份
	var remoteHost = []string{
		confData.REMOTE_HOST,
		confData.REMOTE_PORT,
		confData.REMOTE_USER,
		confData.REMOTE_PASSWORD,
	}
	rmFile := modules.NewBackupClear(confData.SAVE_DAY, remoteHost...)
	rmFile.ClearLocal(confData.BACKUP_SAVE_PATH)
	rmFile.ClearRemote(confData.REMOTE_PATH)

	return nil
}
