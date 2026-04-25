package controller

import (
	"db_backup_go/config"
	"db_backup_go/conn"
	"db_backup_go/logging"
	"db_backup_go/shell"
	"fmt"
	"sync"

	"golang.org/x/crypto/ssh"
)

// 备份主程序，返回 error
func Controller(conf string, dbs string) error {
	//获取配置文件
	execConfig := config.NewConfig(conf)
	err := execConfig.Read()
	if err != nil {
		return err
	}
	//获取要使用的数据库列表
	dbsTxt := config.NewDBList(dbs)
	dbsTxtData, err := dbsTxt.Read()
	if err != nil {
		return err
	}

	dbConnectInfo := &shell.DBInfo{
		DBType:     execConfig.DATABASETYPE,
		ExecPath:   execConfig.MYSQL_EXEC_PATH,
		DBVersion:  execConfig.DB_Version,
		DBHost:     execConfig.DB_HOST,
		DBPort:     execConfig.DB_PORT,
		DBUser:     execConfig.DB_USER,
		DBPassword: execConfig.DB_PASSWORD,
	}
	//对比出要备份的数据库列表
	dbBackupListReference, err := Comparison(execConfig, dbsTxtData, dbConnectInfo)
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

	var sshClient *ssh.Client
	if execConfig.REMOTE_BACKUP {
		logging.Logger.Println("远程备份功能已开启")
		sshClient, err = conn.CreateSSHSocket(execConfig.REMOTE_HOST, execConfig.REMOTE_PORT, execConfig.REMOTE_USER, execConfig.REMOTE_PASSWORD)
		if err != nil {
			return err
		}
	}

	for _, dbBackupReference := range *dbBackupListReference {
		logging.Logger.Printf("%v备份开始", dbBackupReference)
		backuperObject := NewBackuper(execConfig, dbBackupReference)
		wg.Add(1)
		go func(dbBackupReference string) {
			sendRemoteFailed, err := backuperObject.RunBackup(sshClient, dbConnectInfo, dbBackupReference)
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
	backupCleanerObject := NewBackupCleaner(execConfig.SAVE_DAY, dbBackupListReference, sshClient)

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
