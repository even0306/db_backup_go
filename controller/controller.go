package controller

import (
	"db_backup_go/common"
	"db_backup_go/config"
	"db_backup_go/logging"
	"db_backup_go/modules/clear"
	"db_backup_go/modules/database"
	"db_backup_go/modules/run"
	"sync"
)

type Controller interface {
	Controller() error
}

type fileInfo struct {
	confFile string
	dbsFile  string
}

// 初始化控制器，传入配置文件和数据库列表文件，返回 *fileInfo 结构体实例
func NewController(conf string, dbs string) *fileInfo {
	return &fileInfo{
		confFile: conf,
		dbsFile:  dbs,
	}
}

// 备份主程序，返回 error
func (fi fileInfo) Controller() error {
	//获取配置文件
	conf := config.NewConfig(fi.confFile)
	err := conf.Read()
	if err != nil {
		return err
	}
	//获取要使用的数据库列表
	dbs := config.NewDBList(fi.dbsFile)
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

	var wg sync.WaitGroup
	var sendRemoteFailed bool
	//开始循环备份每个数据库
	var responseChannel = make(chan string)
	var sendRemoteFailedChannel = make(chan bool)
	go func() {
		for {
			name := <-responseChannel
			sendRemoteFailed = <-sendRemoteFailedChannel
			logging.Logger.Printf("%v备份完成", name)
			wg.Done()
		}
	}()

	bk := run.NewBackuper(conf)
	for _, v := range *preDBS {
		logging.Logger.Printf("%v备份开始", v)
		wg.Add(1)
		go func(db string) {
			sendRemoteFailed, err := bk.Run(db)
			if err != nil {
				if sendRemoteFailed {
					logging.Logger.Printf("%v发送到异机失败：%v", db, err)
				} else {
					logging.Logger.Panicf("%v备份失败：%v", db, err)
				}
			}
			responseChannel <- db
			sendRemoteFailedChannel <- sendRemoteFailed
		}(v)
	}
	wg.Wait()

	//按天保留最新7份备份，删除之前的备份
	logging.Logger.Printf("开始清理%v天前的备份", conf.SAVE_DAY)
	sshSocket := common.NewSshSocket(conf.REMOTE_HOST, conf.REMOTE_PORT, conf.REMOTE_USER, conf.REMOTE_PASSWORD)
	rmFile := clear.NewBackupClear(conf.SAVE_DAY, preDBS, *sshSocket)

	logging.Logger.Println("开始清理本地备份")
	err = rmFile.ClearLocal(conf.BACKUP_SAVE_PATH)
	if err != nil {
		return err
	} else {
		logging.Logger.Println("本地备份清理完成")
	}

	if conf.REMOTE_BACKUP && !sendRemoteFailed {
		logging.Logger.Println("开始清理远程备份")
		err = rmFile.ClearRemote(conf.REMOTE_PATH)
		if err != nil {
			return err
		} else {
			logging.Logger.Println("远程备份清理完成")
		}
	}
	return nil
}
