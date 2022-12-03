package controller

import (
	"db_backup_go/common"
	"db_backup_go/modules/clear"
	"db_backup_go/modules/database"
	"db_backup_go/modules/run"
	"log"
	"sync"
)

type Controller interface {
	Controller() error
}

type fileInfo struct {
	confFile     string
	dbsFile      string
	fileNameList []string
}

//初始化控制器，传入配置文件和数据库列表文件，返回 *fileInfo 结构体实例
func NewController(conf string, dbs string) *fileInfo {
	return &fileInfo{
		confFile:     conf,
		dbsFile:      dbs,
		fileNameList: []string{},
	}
}

//备份主程序，返回 error
func (fi fileInfo) Controller() error {
	//获取配置文件
	conf := common.NewConfig(fi.confFile)
	err := conf.Read()
	if err != nil {
		return err
	}
	//获取要使用的数据库列表
	dbs := common.NewDBList(fi.dbsFile)
	dbsData, err := dbs.Read()
	if err != nil {
		return err
	}

	//对比出要备份的数据库列表
	cp := database.NewCompartor(conf, dbsData)
	preDBS, err := cp.Comparison()
	if err != nil {
		return err
	}

	//开始循环备份每个数据库
	var responseChannel = make(chan string)
	go func(fl *[]string) {
		for v := range responseChannel {
			*fl = append(*fl, v)
		}
	}(&fi.fileNameList)

	var wg sync.WaitGroup
	limiter := make(chan bool, 4)
	bk := run.NewBackuper(conf)
	for _, v := range *preDBS {
		log.Printf("%v备份开始", v)
		wg.Add(1)
		limiter <- true
		go func(db string) {
			fileName, err := bk.Run(&db)
			if err != nil {
				log.Panicf("%v备份失败：%v", db, err)
			}
			defer wg.Done()
			responseChannel <- fileName
			<-limiter
		}(v)
	}
	wg.Wait()

	//按天保留最新7份备份，删除之前的备份
	sshSocket := common.NewSshSocket(conf.REMOTE_HOST, conf.REMOTE_PORT, conf.REMOTE_USER, conf.REMOTE_PASSWORD)

	rmFile := clear.NewBackupClear(conf.SAVE_DAY, *sshSocket)
	rmFile.ClearLocal(conf.BACKUP_SAVE_PATH)
	rmFile.ClearRemote(conf.REMOTE_PATH)

	for _, v := range fi.fileNameList {
		log.Printf("%v备份成功", v)
	}
	return nil
}
