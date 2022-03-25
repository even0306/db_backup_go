package core

import (
	"fmt"
	"log"
	"mysql_backup_go/config"
	"mysql_backup_go/connection"
	"mysql_backup_go/utils"
	"os/exec"
	"time"
)

type fileName struct {
	filenameNoDate string
	filename       string
	date           string
}

func Backup() {
	//日志初始化
	l := config.Logger{}
	logfile, logs := l.SetLogConfig("server.log")
	log.SetOutput(logfile)

	//获取配置文件
	conf := config.ConfigFile{}
	confData := conf.Read("config.json")
	//获取要使用的数据库列表
	dbs := config.Databases{}
	dbsData := dbs.Read("dbs.txt")
	for _, v := range dbsData {
		if v == "all" {
			dbsData = nil
			dbsData = append(dbsData, "all")
		}
	}

	//获取所有数据库名
	cmd := exec.Command(confData.MYSQL_EXEC_PATH+"/mysql", "-h", confData.DB_HOST, "-P", string(confData.DB_PORT), "-u", confData.DB_USER, "-p"+confData.DB_PASSWORD, "-Bse", `show databases`)
	out, err := cmd.Output()
	if err != nil {
		logs.InfoLogger.Printf("stderr: %v", err)
	}

	//根据筛选方式，筛选出待备份的数据库
	var preDBS []string
	var b = 0
	if confData.FILTER_METHOD == true {
		for _, v := range string(out[:]) {
			for _, w := range dbsData {
				if w == "all" {
					preDBS = nil
					preDBS = append(preDBS, "all")
					b = 1
					break
				} else if string(v) == w {
					preDBS = append(preDBS, w)
				}
			}
			if b == 1 {
				break
			}
		}
	} else {
		for _, v := range string(out[:]) {
			for _, w := range dbsData {
				if w == "all" {
					preDBS = nil
					preDBS = append(preDBS, "all")
					b = 1
					break
				} else if string(v) != w {
					preDBS = append(preDBS, w)
				}
			}
			if b == 1 {
				break
			}
		}
	}

	//开始循环备份每个数据库
	var fileName fileName
	for _, v := range preDBS {
		fileName.date = time.Now().Format("2020-01-02 15:04:02")
		fileName.filenameNoDate = v + "_" + confData.DB_LABEL
		fileName.filename = fileName.filenameNoDate + "_" + fileName.date
		fmt.Printf("do backup %v", v)
		if v == "all" {
			cmd := exec.Command(confData.MYSQL_EXEC_PATH+"/mysqldump", "-h", confData.DB_HOST, "-P", string(confData.DB_PORT), "-u", confData.DB_USER, "-p"+confData.DB_PASSWORD, "-E", "-R", "--triggers", "--all-databases")
			out, err = cmd.Output()
		} else {
			cmd := exec.Command(confData.MYSQL_EXEC_PATH+"/mysqldump", "-h", confData.DB_HOST, "-P", string(confData.DB_PORT), "-u", confData.DB_USER, "-p"+confData.DB_PASSWORD, "-E", "-R", "--triggers", v)
			out, err = cmd.Output()
		}
		if err != nil {
			logs.ErrorLogger.Panicf(v+"数据库配置失败：%v", err)
		}

		//压缩并保存备份文件
		saveFile := Gz{}
		saveFile.CompressFile(out, confData.BACKUP_SAVE_PATH, fileName.filename)
		if err != nil {
			logs.ErrorLogger.Panicf("文件保存失败：%v", err)
		}

		// 判断是否开启远程备份功能
		if confData.REMOTE_BACKUP == true {
			//判断远端系统类型
			tp := utils.Type{}
			//根据远端系统类型，发送备份文件到远端
			if tp.CheckOS() == "linux" {
				sendToLinux := connection.Linux{}
				sendToLinux.SendToRemoteHost()
			} else if tp.CheckOS() == "windows" {
				sendToWindows := connection.Windows{}
				sendToWindows.SendToRemoteHost()
			}
		}
	}
}
