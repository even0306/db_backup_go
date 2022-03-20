package core

import (
	"io"
	"mysql_backup_go/config"
	"os/exec"
)

type fileName struct {
	filenameNoDate string
	filename       string
	date           string
}

func Backup() error {
	conf := config.ConfigFile{}
	confData, err := conf.Read("config.json")
	if err != nil {
		return err
	}
	dbs := config.Databases{}
	dbsData, err := dbs.Read("dbs.txt")
	if err != nil && err != io.EOF {
		return err
	}
	for _, v := range dbsData {
		if v == "all" {
			dbsData = nil
			dbsData = append(dbsData, "all")
		}
	}

	cmd := exec.Command(confData.MYSQL_EXEC_PATH+"/mysql", "-h", confData.DB_HOST, "-P", string(confData.DB_PORT), "-u", confData.DB_USER, "-p"+confData.DB_PASSWORD, "-Bse", `show databases`)
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	var preDBS []string
	var b = 0
	if confData.FILTER_METHOD == true {
		for _, v := range string(out[:]) {
			for _, w := range dbsData {
				if w == "all" {
					preDBS = nil
					preDBS = append(preDBS, "--all-databases")
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
					preDBS = append(preDBS, "--all-databases")
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

	for _, v := range preDBS {

	}

	return err
}
