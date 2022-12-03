package run

import (
	"db_backup_go/common"
	"db_backup_go/modules"
	"db_backup_go/modules/dump"
	"db_backup_go/modules/ssh"
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

	var out *[]byte
	var err error
	dbi := dump.DBInfo{
		DBHost:     b.conf.DB_HOST,
		DBPort:     b.conf.DB_PORT,
		DBUser:     b.conf.DB_USER,
		DBPassword: b.conf.DB_PASSWORD,
	}
	dbu := dump.NewDBDumpFunc(b.conf.MYSQL_EXEC_PATH, &dbi)
	if b.conf.DATABASETYPE == "mysql" {
		if *db == "all" {
			out, err = dbu.MysqlDumpAll()
			if err != nil {
				return "", err
			}
		} else {
			out, err = dbu.MysqlDump(db)
			if err != nil {
				return "", err
			}
		}
	} else if b.conf.DATABASETYPE == "postgresql" {
		if *db == "all" {
			out, err = dbu.PostgresqlDumpAll()
			if err != nil {
				return "", err
			}
		} else {
			out, err = dbu.PostgresqlDump(db)
			if err != nil {
				return "", err
			}
		}
	}

	//压缩并保存备份文件
	saveFile := modules.NewCompress(out, &fileName)
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
	if b.conf.REMOTE_BACKUP == true {
		//发送备份文件到远端
		s := modules.NewSshSocket(b.conf.REMOTE_HOST, b.conf.REMOTE_PORT, b.conf.REMOTE_USER, b.conf.REMOTE_PASSWORD)
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

		up := ssh.NewSftpOperater(sftpClient)
		err = up.Upload(b.conf.BACKUP_SAVE_PATH+"/"+*db+"/"+fileName+".gz", b.conf.REMOTE_PATH+"/"+*db, fileName+".gz")
		if err != nil {
			return "", err
		}
	}

	return *db, nil
}
