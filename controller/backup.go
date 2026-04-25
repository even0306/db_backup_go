package controller

import (
	"db_backup_go/config"
	"db_backup_go/conn"
	"db_backup_go/shell"
	"fmt"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type BackupInfo struct {
	Conf *config.ConfigFile

	Date string

	MyBackupPath config.BackupPath
}

// 初始化备份工具，传入*common.ConfigFile类型的配置数据
func NewBackuper(conf *config.ConfigFile, db string) *BackupInfo {
	date := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("%v_%v_%v.sql.gz", db, conf.DB_LABEL, date)
	return &BackupInfo{
		Conf: conf,
		Date: date,
		MyBackupPath: config.BackupPath{
			FileName:               fileName,
			SavePath:               fmt.Sprintf("%v/%v", conf.BACKUP_SAVE_PATH, conf.DB_LABEL),
			FullSavePath:           fmt.Sprintf("%v/%v/%v", conf.BACKUP_SAVE_PATH, conf.DB_LABEL, db),
			FullSavePathFile:       fmt.Sprintf("%v/%v/%v/%v", conf.BACKUP_SAVE_PATH, conf.DB_LABEL, db, fileName),
			RemoteSavePath:         fmt.Sprintf("%v/%v", conf.REMOTE_PATH, conf.DB_LABEL),
			RemoteFullSavePath:     fmt.Sprintf("%v/%v/%v", conf.REMOTE_BACKUP, conf.DB_LABEL, db),
			RemoteFullSavePathFile: fmt.Sprintf("%v/%v/%v/%v", conf.REMOTE_BACKUP, conf.DB_LABEL, db, fileName),
		},
	}
}

// 执行备份，传入ssh客户端，返回是否发送远程失败和error
func (backupInfo *BackupInfo) RunBackup(sshClient *ssh.Client, dbConnectInfo *shell.DBInfo, dbBackupReference string) (bool, error) {
	//根据数据库类型选择相应的备份工具,传入备份工具需要的参数，执行备份命令，返回备份结果和错误
	err := shell.BackupSelecter(dbConnectInfo, dbBackupReference, backupInfo.MyBackupPath, backupInfo.Conf.SINGLE_TRANSACTION)
	if err != nil {
		return false, err
	}

	// 判断是否开启远程备份功能
	if sshClient != nil {
		sftpClient, err := sftp.NewClient(sshClient)
		if err != nil {
			return true, fmt.Errorf("创建sftp客户端失败：%w", err)
		}
		defer sftpClient.Close()

		up := conn.NewSftpOperater(sftpClient)
		err = up.Upload(backupInfo.MyBackupPath.FullSavePathFile, fmt.Sprintf("%v/%v/%v", backupInfo.Conf.REMOTE_PATH, backupInfo.Conf.DB_LABEL, dbBackupReference), backupInfo.MyBackupPath.FileName)
		if err != nil {
			return true, err
		}
	}

	return false, nil
}
