package modules

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mysql_backup_go/common"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/sftp"
)

type Backup interface {
	Run(db *string) (string, error)
}

type backupInfo struct {
	conf *common.ConfigFile

	fileNameNoDate string
	fileName       string
	date           string
}

func NewBackuper(conf *common.ConfigFile) *backupInfo {
	return &backupInfo{
		conf:           conf,
		fileNameNoDate: "",
		fileName:       "",
		date:           "",
	}
}

//循环备份每个数据库，返回本地备份位置，异机备份位置，和err
func (b *backupInfo) Run(db *string) (string, error) {
	b.date = time.Now().Format("2006-01-02")
	b.fileNameNoDate = *db + "_" + b.conf.DB_LABEL
	b.fileName = b.fileNameNoDate + "_" + b.date + ".sql"
	log.Printf("正在备份：%v", *db)
	var out []byte
	var err error
	if *db == "all" {
		cmd := exec.Command(b.conf.MYSQL_EXEC_PATH+"/mysqldump", "-h"+b.conf.DB_HOST, "-P"+string(b.conf.DB_PORT), "-u"+b.conf.DB_USER, "-p"+b.conf.DB_PASSWORD, "-E", "-R", "--triggers", "--all-databases")
		out, err = cmd.Output()
	} else {
		cmd := exec.Command(b.conf.MYSQL_EXEC_PATH+"/mysqldump", "-h"+b.conf.DB_HOST, "-P"+string(b.conf.DB_PORT), "-u"+b.conf.DB_USER, "-p"+b.conf.DB_PASSWORD, "-E", "-R", "--triggers", *db)
		out, err = cmd.Output()
	}
	if err != nil {
		return "", fmt.Errorf(*db+"数据库配置失败：%w", err)
	}

	//压缩并保存备份文件
	saveFile := NewCompress(&out, &b.fileName)
	buff, err := saveFile.CompressFile()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(b.conf.BACKUP_SAVE_PATH+"/"+*db, 0777)
	if err != nil {
		return "", fmt.Errorf("创建备份文件路径失败：%w", err)
	}
	f, err := os.Create(b.conf.BACKUP_SAVE_PATH + "/" + *db + "/" + b.fileName + ".gz")
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
		s := NewSshSocket(b.conf.REMOTE_HOST, b.conf.REMOTE_PORT, b.conf.REMOTE_USER, b.conf.REMOTE_PASSWORD)
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

		up := NewSftpOperater(sftpClient)
		err = up.Upload(b.conf.BACKUP_SAVE_PATH+b.fileName, b.conf.REMOTE_PATH)
		if err != nil {
			return "", err
		}
	}

	return b.fileName, nil
}
