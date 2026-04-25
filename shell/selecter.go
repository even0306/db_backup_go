package shell

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"db_backup_go/config"
	"db_backup_go/logging"
	"fmt"
	"io"
	"os"
)

type DBInfo struct {
	DBType     string
	ExecPath   string
	DBVersion  string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
}

// 备份工具选择器，传入 *common.ConfigFile 和要备份的库名指针，返回备份出的字节流指针和报错信息
func BackupSelecter(dbInfo *DBInfo, dbBackupReference string, backupPath config.BackupPath, single int) error {
	var err error
	switch dbInfo.DBType {
	case "mysql", "mariadb":
		if dbBackupReference == "all" {
			backupPath.FullSavePath = fmt.Sprintf("%v/all", backupPath.SavePath)
			backupPath.FullSavePathFile = fmt.Sprintf("%v/all/%v", backupPath.SavePath, backupPath.FileName)
			err = MysqlDumpAll(dbInfo, backupPath, single)
			if err != nil {
				return err
			}
		} else {
			err = MysqlDump(dbInfo, backupPath, single, dbBackupReference)
			if err != nil {
				return err
			}
		}
	case "postgresql":
		if dbBackupReference == "all" {
			backupPath.FullSavePath = fmt.Sprintf("%v/all", backupPath.SavePath)
			backupPath.FullSavePathFile = fmt.Sprintf("%v/all/%v", backupPath.SavePath, backupPath.FileName)
			err = PostgresqlDumpAll(dbInfo, backupPath)
			if err != nil {
				return err
			}
		} else {
			err = PostgresqlDump(dbInfo, backupPath, dbBackupReference)
			if err != nil {
				return err
			}
		}
	default:
		logging.Logger.Panic("未知的数据库类型，请重新检查 config.json 文件配置")
	}
	return nil
}

func DBListSelecter(dbInfo *DBInfo) (*[]string, error) {
	switch dbInfo.DBType {
	case "mysql", "mariadb":
		allDbs, err := GetMysqlDBList(dbInfo)
		if err != nil {
			return nil, err
		}
		return allDbs, nil
	case "postgresql":
		allDbs, err := GetPostgresqlDBList(dbInfo)
		if err != nil {
			return nil, err
		}
		return allDbs, nil
	default:
		logging.Logger.Panic("未知的数据库类型，请重新检查 config.json 文件配置")
	}
	return nil, nil
}

func WriteToFile(stdout *io.ReadCloser, stderr *bytes.Buffer, savePath config.BackupPath) {
	var gz *gzip.Writer
	reader := bufio.NewReader(*stdout)
	isFirst := true
	for {
		line, err := reader.ReadString('\n')
		if isFirst {
			if line == "" {
				logging.Logger.Panic(stderr.String())
			} else {
				err = os.MkdirAll(savePath.FullSavePath, 0777)
				if err != nil {
					logging.Logger.Panic(err)
				}

				f, err := os.Create(savePath.FullSavePathFile)
				if err != nil {
					logging.Logger.Panic(err)
				}

				//创建一个gzip的流来接收管道中内容
				gz = gzip.NewWriter(f)
				defer gz.Close()
			}
			isFirst = false
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			logging.Logger.Panicf("读取流出现问题：%v，文件备份不完整。", err)
		}
		_, err = gz.Write([]byte(line))
		if err != nil {
			logging.Logger.Panic(err)
		}
	}
}
