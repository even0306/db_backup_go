package controller

import (
	"db_backup_go/config"
	"db_backup_go/shell"
	"fmt"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type backupInfo struct {
	conf *config.ConfigFile

	date string
}

// 初始化备份工具，传入*common.ConfigFile类型的配置数据
func NewBackuper(conf *config.ConfigFile) *backupInfo {
	return &backupInfo{
		conf: conf,
		date: "",
	}
}

// 循环备份每个数据库，返回库名或err
func (b *backupInfo) RunBackup(db string, sshClient *ssh.Client) (bool, error) {
	b.date = time.Now().Format("2006-01-02")
	fileName := db + "_" + b.conf.DB_LABEL + "_" + b.date + ".sql.gz"

	//根据数据库类型选择相应的备份工具
	dbi := shell.NewSelecter(b.conf.DATABASETYPE, b.conf.MYSQL_EXEC_PATH, b.conf.DB_Version, b.conf.DB_HOST, b.conf.DB_PORT, b.conf.DB_USER, b.conf.DB_PASSWORD)
	err := shell.BackupSelecter(dbi, db, b.conf.BACKUP_SAVE_PATH+"/"+b.conf.DB_LABEL, fileName, b.conf.SINGLE_TRANSACTION)
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

		up := NewSftpOperater(sftpClient)
		err = up.Upload(fmt.Sprintf("%v/%v/%v/%v", b.conf.BACKUP_SAVE_PATH, b.conf.DB_LABEL, db, fileName), fmt.Sprintf("%v/%v/%v", b.conf.REMOTE_PATH, b.conf.DB_LABEL, db), fileName)
		if err != nil {
			return true, err
		}
	}

	return false, nil
}
