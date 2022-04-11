package modules

import (
	"db_backup_go/common"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/sftp"
)

type Clear interface {
	ClearLocal(dict string) error
	ClearRemote(dict string) error
}

type remoteHost struct {
	host     string
	port     int
	user     string
	password string
}

type backupFile struct {
	ConnInfo
	saveDay int
}

func NewBackupClear(saveDay int, ci ConnInfo) *backupFile {
	return &backupFile{
		ConnInfo: ci,
		saveDay:  saveDay,
	}
}

func (bf *backupFile) ClearLocal(dict string) error {
	//确认要保留的文件
	fsDict, err := ioutil.ReadDir(dict)
	if err != nil {
		return fmt.Errorf("读取目录失败：%w", err)
	}
	cf := common.NewOrder(fsDict)
	sfDict := cf.SortByTime()
	if len(sfDict) < bf.saveDay {
		bf.saveDay = len(sfDict)
	}

	sfDict = sfDict[bf.saveDay:]

	//删除旧备份
	for _, v := range sfDict {
		os.Remove(dict + "/" + v.Name())
	}

	return nil
}

func (bf *backupFile) ClearRemote(dict string) error {
	//确认要保留的文件
	sshSocket := NewSshSocket(bf.ConnInfo.Host, bf.ConnInfo.Port, bf.ConnInfo.User, bf.ConnInfo.Password)
	sshClient, err := sshSocket.Connect()
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	fsDict, err := sftpClient.ReadDir(dict)
	if err != nil {
		return fmt.Errorf("读取远程目录失败：%w", err)
	}
	cf := common.NewOrder(fsDict)
	sfDict := cf.SortByTime()
	if len(sfDict) < bf.saveDay {
		bf.saveDay = len(sfDict)
	}

	sfDict = sfDict[bf.saveDay:]

	//删除旧备份
	cmd := NewSftpOperater(sftpClient)
	for _, v := range sfDict {
		err := cmd.Remove(dict + "/" + v.Name())
		if err != nil {
			return fmt.Errorf("删除远程目录文件失败：%w", err)
		}
	}

	return nil
}
