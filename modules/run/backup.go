package run

import (
	"db_backup_go/common"
	"db_backup_go/modules/send"
	"db_backup_go/shell"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
)

type Backup interface {
	Run(db *string) (string, error)
}

type backupInfo struct {
	conf *common.ConfigFile

	date string
}

//初始化备份工具，传入*common.ConfigFile类型的配置数据
func NewBackuper(conf *common.ConfigFile) *backupInfo {
	return &backupInfo{
		conf: conf,
		date: "",
	}
}

//循环备份每个数据库，返回本地备份位置，异机备份位置，和err
func (b *backupInfo) Run(db *string) (string, error) {
	b.date = time.Now().Format("2006-01-02")
	fileNameNoDate := *db + "_" + b.conf.DB_LABEL
	fileName := fileNameNoDate + "_" + b.date + ".sql"

	//根据数据库类型选择相应的备份工具
	dbi := shell.NewSelecter(b.conf.DATABASETYPE, b.conf.MYSQL_EXEC_PATH, b.conf.DB_Version, b.conf.DB_HOST, b.conf.DB_PORT, b.conf.DB_USER, b.conf.DB_PASSWORD)
	out, err := shell.BackupSelecter(dbi, db)
	if err != nil {
		return "", err
	}

	//压缩并保存备份文件
	saveFile := common.NewCompress(out, &fileName)
	buff, err := saveFile.CompressFile()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(b.conf.BACKUP_SAVE_PATH+"/"+*db, 0777)
	if err != nil {
		return "", fmt.Errorf("创建备份文件路径失败：%w", err)
	}
	f, err := os.Create(b.conf.BACKUP_SAVE_PATH + "/" + *db + "/" + fileName + ".gz")
	if err != nil {
		return "", fmt.Errorf("创建备份文件失败：%w", err)
	}

	defer func() error {
		err = f.Close()
		if err != nil {
			return fmt.Errorf("压缩数据传输进os.file最终关闭失败：%w", err)
		}
		return nil
	}()

	_, err = io.Copy(f, buff)
	if err != nil {
		errors.Unwrap(err)
		return "", fmt.Errorf("压缩数据传输进os.file失败：%w", err)
	}

	// 判断是否开启远程备份功能
	if b.conf.REMOTE_BACKUP {
		//发送备份文件到远端
		s := common.NewSshSocket(b.conf.REMOTE_HOST, b.conf.REMOTE_PORT, b.conf.REMOTE_USER, b.conf.REMOTE_PASSWORD)
		sshClient, err := s.Connect()
		if err != nil {
			return "", err
		}
		defer sshClient.Close()

		sftpClient, err := sftp.NewClient(sshClient)
		if err != nil {
			return "", fmt.Errorf("创建sftp客户端失败：%w", err)
		}
		defer sftpClient.Close()

		up := send.NewSftpOperater(sftpClient)
		err = up.Upload(b.conf.BACKUP_SAVE_PATH+"/"+*db+"/"+fileName+".gz", b.conf.REMOTE_PATH+"/"+*db, fileName+".gz")
		if err != nil {
			return "", err
		}
	}

	return *db, nil
}
