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

type backupFile struct {
	ConnInfo
	saveDay int
}

//初始化旧备份清理，传入保存的天数和远端服务器连接信息（ConnInfo结构体）
func NewBackupClear(saveDay int, sc ConnInfo) *backupFile {
	return &backupFile{
		ConnInfo: sc,
		saveDay:  saveDay,
	}
}

//清理本地旧备份文件，传入本地路径，返回error
func (bf *backupFile) ClearLocal(dict string) error {
	//确认要保留的文件
	fsDict, err := ioutil.ReadDir(dict)
	if err != nil {
		return fmt.Errorf("读取目录失败：%w", err)
	}
	cf := common.SortByTime(fsDict)
	if len(cf) < bf.saveDay {
		bf.saveDay = len(cf)
	}

	cf = cf[bf.saveDay:]

	//删除旧备份
	for _, v := range cf {
		os.Remove(dict + "/" + v.Name())
	}

	return nil
}

//清理远端旧备份文件，传入远端机器路径，返回error
func (bf *backupFile) ClearRemote(dict string) error {
	//确认要保留的文件
	sshClient, err := bf.Connect()
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
	cf := common.SortByTime(fsDict)
	if len(cf) < bf.saveDay {
		bf.saveDay = len(cf)
	}

	cf = cf[bf.saveDay:]

	//删除旧备份
	cmd := NewSftpOperater(sftpClient)
	for _, v := range cf {
		err := cmd.Remove(dict + "/" + v.Name())
		if err != nil {
			return fmt.Errorf("删除远程目录文件失败：%w", err)
		}
	}

	return nil
}
