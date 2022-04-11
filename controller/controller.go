package controller

import (
	"db_backup_go/common"
	"db_backup_go/modules"
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

func NewController(conf string, dbs string) *fileInfo {
	return &fileInfo{
		confFile:     conf,
		dbsFile:      dbs,
		fileNameList: []string{},
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
	log.Println("备份开始")
	var responseChannel = make(chan string)
	go func(fl *[]string) {
		for v := range responseChannel {
			*fl = append(*fl, v)
		}
	}(&fi.fileNameList)

	var wg sync.WaitGroup
	limiter := make(chan bool, 4)
	bk := modules.NewBackuper(confData)
	for _, v := range *preDBS {
		wg.Add(1)
		limiter <- true
		go func(db *string) {
			fileName, err := bk.Run(db)
			if err != nil {
				log.Printf("%v备份失败：%v", *db, err)
			}
			defer wg.Done()
			log.Printf("%v备份完成", *db)
			responseChannel <- fileName
			<-limiter
		}(&v)
	}
	wg.Wait()

	//按天保留最新7份备份，删除之前的备份
	rh := modules.ConnInfo{
		Host:     confData.REMOTE_HOST,
		Port:     confData.REMOTE_PORT,
		User:     confData.REMOTE_USER,
		Password: confData.REMOTE_PASSWORD,
	}

	rmFile := modules.NewBackupClear(confData.SAVE_DAY, rh)
	rmFile.ClearLocal(confData.BACKUP_SAVE_PATH)
	rmFile.ClearRemote(confData.REMOTE_PATH)

	for _, v := range fi.fileNameList {
		log.Printf("备份结束\n本地备份路径：%v\n远程备份路径(如开启远程备份)：%v", confData.BACKUP_SAVE_PATH+v, confData.REMOTE_PATH+v)
	}
	return nil
}
