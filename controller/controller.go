package controller

import (
	"db_backup_go/common"
	"db_backup_go/config"
	"db_backup_go/logging"
	"db_backup_go/modules/clear"
	"db_backup_go/modules/database"
	"db_backup_go/modules/run"
	"fmt"
	"sync"
)

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
	execConfig := config.NewConfig(fi.confFile)
	err := execConfig.Read()
	if err != nil {
		return err
	}
	//获取要使用的数据库列表
	dbsTxt := config.NewDBList(fi.dbsFile)
	dbsTxtData, err := dbsTxt.Read()
	if err != nil {
		return err
	}

	//对比出要备份的数据库列表
	compartorObject := database.NewCompartor(execConfig, dbsTxtData)
	dbBackupListReference, err := compartorObject.Comparison()
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

	backuperObject := run.NewBackuper(execConfig)
	for _, dbBackupReference := range *dbBackupListReference {
		logging.Logger.Printf("%v备份开始", dbBackupReference)
		wg.Add(1)
		go func(dbBackupReference string) {
			sendRemoteFailed, err := backuperObject.Run(dbBackupReference)
			if err != nil {
				if sendRemoteFailed {
					logging.Logger.Printf("%v发送到异机失败：%v", dbBackupReference, err)
				} else {
					logging.Logger.Panicf("%v备份失败：%v", dbBackupReference, err)
				}
			}
			responseChannel <- dbBackupReference
			sendRemoteFailedChannel <- sendRemoteFailed
		}(dbBackupReference)
	}
	wg.Wait()

	//按天保留最新7份备份，删除之前的备份
	logging.Logger.Printf("开始清理%v天前的备份", execConfig.SAVE_DAY)
	sshSocketCreaterObject := common.NewSSHSocketCreater(execConfig.REMOTE_HOST, execConfig.REMOTE_PORT, execConfig.REMOTE_USER, execConfig.REMOTE_PASSWORD)
	backupCleanerObject := clear.NewBackupCleaner(execConfig.SAVE_DAY, dbBackupListReference, *sshSocketCreaterObject)

	logging.Logger.Println("开始清理本地备份")
	deadFileNameList, err := backupCleanerObject.ClearLocal(fmt.Sprintf("%v/%v", execConfig.BACKUP_SAVE_PATH, execConfig.DB_LABEL))
	if err != nil {
		return err
	} else {
		logging.Logger.Println("本地备份清理完成")
	}

	if execConfig.REMOTE_BACKUP && !sendRemoteFailed {
		logging.Logger.Println("开始清理远程备份")
		err = backupCleanerObject.ClearRemote(fmt.Sprintf("%v/%v", execConfig.REMOTE_PATH, execConfig.DB_LABEL), deadFileNameList)
		if err != nil {
			return err
		} else {
			logging.Logger.Println("远程备份清理完成")
		}
	}
	return nil
}
