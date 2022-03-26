package core

import (
	"fmt"
	"io"
	"log"
	"mysql_backup_go/config"
	"mysql_backup_go/connection"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/sftp"
)

type backup interface {
	Run() error
}

type Backup struct {
	conf *config.ConfigFile
	db   *string

	filenameNoDate string
	filename       string
	date           string
}

func NewBackuper(conf *config.ConfigFile, db *string) *Backup {
	return &Backup{
		conf:           conf,
		db:             db,
		filenameNoDate: "",
		filename:       "",
		date:           "",
	}
}

//循环备份每个数据库
func (b *Backup) Run() error {
	b.date = time.Now().Format("2020-01-02 15:04")
	b.filenameNoDate = *b.db + "_" + b.conf.DB_LABEL
	b.filename = b.filenameNoDate + "_" + b.date + ".sql"
	fmt.Printf("do backup %v", b.db)
	var out []byte
	var err error
	if *b.db == "all" {
		cmd := exec.Command(b.conf.MYSQL_EXEC_PATH+"/mysqldump", "-h"+b.conf.DB_HOST, "-P"+string(b.conf.DB_PORT), "-u"+b.conf.DB_USER, "-p"+b.conf.DB_PASSWORD, "-E", "-R", "--triggers", "--all-databases")
		out, err = cmd.Output()
	} else {
		cmd := exec.Command(b.conf.MYSQL_EXEC_PATH+"/mysqldump", "-h"+b.conf.DB_HOST, "-P"+string(b.conf.DB_PORT), "-u"+b.conf.DB_USER, "-p"+b.conf.DB_PASSWORD, "-E", "-R", "--triggers", *b.db)
		out, err = cmd.Output()
	}
	if err != nil {
		return fmt.Errorf(*b.db+"数据库配置失败：%v", err)
	}

	//压缩并保存备份文件
	saveFile := NewCompress(&out, &b.filename)
	buff, err := saveFile.CompressFile()
	if err != nil {
		return err
	}

	f, err := os.Create(b.conf.BACKUP_SAVE_PATH + b.filename + ".gz")
	if err != nil {
		return fmt.Errorf("创建备份文件失败：%v", err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Panicf("压缩数据传输进os.file最终关闭失败：%v", err)
		}
	}()

	_, err = io.Copy(f, buff)
	if err != nil {
		return fmt.Errorf("压缩数据传输进os.file失败：%v", err)
	}

	// 判断是否开启远程备份功能
	if b.conf.REMOTE_BACKUP == true {
		//发送备份文件到远端
		s := connection.NewSshSocket(b.conf)
		sshClient, err := s.Connect()
		if err != nil {
			return err
		}
		sftpClient, err := sftp.NewClient(sshClient)
		if err != nil {
			return fmt.Errorf("创建sftp客户端失败：%v", err)
		}

		up := NewUploader(sftpClient, b.conf.BACKUP_SAVE_PATH+b.filename, b.conf.REMOTE_PATH)
		err = up.Upload()
		if err != nil {
			return err
		}
	}

	return nil
}
