package core

import (
	"fmt"
	"log"
	"mysql_backup_go/config"
	"os/exec"
	"time"
)

type fileName struct {
	filenameNoDate string
	filename       string
	date           string
}

func Backup() {
	l := config.Logger{}
	logfile, logs := l.SetLogConfig("server.log")
	log.SetOutput(logfile)
	conf := config.ConfigFile{}
	confData := conf.Read("config.json")
	dbs := config.Databases{}
	dbsData := dbs.Read("dbs.txt")
	for _, v := range dbsData {
		if v == "all" {
			dbsData = nil
			dbsData = append(dbsData, "all")
		}
	}

	cmd := exec.Command(confData.MYSQL_EXEC_PATH+"/mysql", "-h", confData.DB_HOST, "-P", string(confData.DB_PORT), "-u", confData.DB_USER, "-p"+confData.DB_PASSWORD, "-Bse", `show databases`)
	out, err := cmd.Output()
	if err != nil {
		logs.InfoLogger.Printf("stderr: %v", err)
	}

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

		saveFile := File{}
		saveFile.SaveFile(out, confData.BACKUP_SAVE_PATH, fileName.filename)
		if err != nil {
			logs.ErrorLogger.Panicf("文件保存失败：%v", err)
		}
	}
}
