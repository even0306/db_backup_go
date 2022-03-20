package core

import (
	"fmt"
	"mysql_backup_go/config"
	"os/exec"
)

func Backup() error {
	conf := config.Config{}
	confData, err := conf.Read("config.json")
	if err != nil {
		return err
	}
	dbs := config.Databases{}
	dbsData, err := dbs.Read("dbs.txt")
	for _, v := range dbsData {
		if v == "all" {
			dbsData = nil
			dbsData = append(dbsData, "all")
		}
	}

	if conf.FILTER_METHOD == true {
		cmd := exec.Command(confData.MYSQL_EXEC_PATH+"/mysql", "-h"+confData.DB_HOST+" -P"+confData.DB_PORT+" -u"+confData.DB_USER+" -p"+confData.DB_PASSWORD+" -Bse 'show databases' | grep -f ./dbs.txt")
		fmt.Println(cmd)
	} else {
		cmd := exec.Command(confData.MYSQL_EXEC_PATH+"/mysql", "-h"+confData.DB_HOST+" -P"+confData.DB_PORT+" -u"+confData.DB_USER+" -p"+confData.DB_PASSWORD+" -Bse 'show databases' | grep -v -f ./dbs.txt")
		fmt.Println(cmd)
	}

	return err
}
