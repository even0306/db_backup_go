package clear

import (
	"db_backup_go/common"
	"db_backup_go/logging"
	"db_backup_go/modules/send"
	"fmt"
	"io/fs"
	"os"

	"github.com/pkg/sftp"
)

type databaseBackuper struct {
	common.ConnInfo
	saveDay               int
	dbBackupListReference *[]string
}

// 初始化旧备份清理，传入保存的天数和远端服务器连接信息（ConnInfo结构体）
func NewBackupCleaner(saveDay int, dbBackupListReference *[]string, sshSocketCreaterObject common.ConnInfo) *databaseBackuper {
	return &databaseBackuper{
		ConnInfo:              sshSocketCreaterObject,
		saveDay:               saveDay,
		dbBackupListReference: dbBackupListReference,
	}
}

// 清理本地旧备份文件，传入本地路径，返回error
func (bker *databaseBackuper) ClearLocal(backupSavePath string) ([]string, error) {
	//确认要保留的文件
	backupSavePathObjects, err := os.ReadDir(backupSavePath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败：%w", err)
	}
	var backupSavePathFileNameList []string
	for _, backupSavePathFile := range backupSavePathObjects {
		if backupSavePathFile.IsDir() {
			backupSavePathFileNameList = append(backupSavePathFileNameList, backupSavePathFile.Name())
		}
	}

	var deadFileNameList []string
	var backupSavePathFileObject []fs.DirEntry
	for _, backupSavePathFileName := range backupSavePathFileNameList {
		execStop := false
		for i, dbBackupReference := range *bker.dbBackupListReference {
			if i > len(*bker.dbBackupListReference) || backupSavePathFileName == dbBackupReference {
				execStop = false
				break
			}
			execStop = true
		}

		if execStop {
			continue
		}

		backupSavePathFileObject, err = os.ReadDir(backupSavePath + "/" + backupSavePathFileName)
		if err != nil {
			return nil, fmt.Errorf("读取目录下文件失败：%w", err)
		}

		backupSavePathFileListDESC := common.SortByTime(backupSavePathFileObject)

		deadDay := bker.saveDay
		if len(backupSavePathFileListDESC) < bker.saveDay {
			deadDay = len(backupSavePathFileListDESC)
		}

		//
		emptyFileNum := 0
		for index, backupSavePathFile := range backupSavePathFileListDESC {
			if index == deadDay {
				break
			}

			backupSavePathFileByte, err := os.ReadFile(backupSavePath + "/" + backupSavePathFileName + "/" + backupSavePathFile.Name())
			if err != nil {
				return nil, err
			}

			if len(backupSavePathFileByte) < 400 {
				emptyFileNum += 1
			}
		}

		deadDay = deadDay + emptyFileNum

		var backupSavePathDeadFileList []fs.DirEntry
		if len(backupSavePathFileListDESC) < deadDay {
			backupSavePathFileListDESC = nil
		} else {
			backupSavePathDeadFileList = backupSavePathFileListDESC[deadDay:]
		}

		//删除旧备份
		for _, deadFile := range backupSavePathDeadFileList {
			deadFileNameList = append(deadFileNameList, deadFile.Name())
			logging.Logger.Printf("删除本地文件: %v/%v/%v", backupSavePath, backupSavePathFileName, deadFile.Name())
			err := os.Remove(backupSavePath + "/" + backupSavePathFileName + "/" + deadFile.Name())
			if err != nil {
				return nil, fmt.Errorf("旧备份文件删除失败：%w", err)
			}
		}

		//检查是否还存在指定份数的备份
		backupSavePathFileNameList, err := os.ReadDir(backupSavePath + "/" + backupSavePathFileName)
		if err != nil {
			return nil, fmt.Errorf("读取目录失败：%w", err)
		}
		if len(backupSavePathFileNameList)-emptyFileNum < bker.saveDay {
			logging.Logger.Printf("%v有效备份数：%v,不足%v份", backupSavePathFileName, len(backupSavePathFileNameList)-emptyFileNum, bker.saveDay)
		}
	}

	return deadFileNameList, nil
}

// 清理远端旧备份文件，传入远端机器路径和要删除的文件名列表，返回error
func (bker *databaseBackuper) ClearRemote(backupSavePath string, deadFileNameList []string) error {
	//确认要保留的文件
	sshClient, err := bker.Connect()
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	backupSavePathName, err := sftpClient.ReadDir(backupSavePath)
	if err != nil {
		return err
	}

	//删除旧备份
	cmd := send.NewSftpOperater(sftpClient)
	for _, deadFileName := range deadFileNameList {
		logging.Logger.Printf("删除远程文件: %v/%v/%v", backupSavePath, backupSavePathName, deadFileName)
		err := cmd.Remove(fmt.Sprintf("%v/%v/%v", backupSavePath, backupSavePathName, deadFileName))
		if err != nil {
			return fmt.Errorf("删除远程目录文件失败：%w", err)
		}
	}

	return nil
}
