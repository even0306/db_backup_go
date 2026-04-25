package controller

import (
	"db_backup_go/config"
	"db_backup_go/logging"
	"db_backup_go/shell"
	"strings"
)

func Comparison(conf *config.ConfigFile, dbs *[]string) (*[]string, error) {
	//获取所有数据库名
	dbi := shell.NewSelecter(conf.DATABASETYPE, conf.MYSQL_EXEC_PATH, conf.DB_Version, conf.DB_HOST, conf.DB_PORT, conf.DB_USER, conf.DB_PASSWORD)
	allDbs, err := shell.DBListSelecter(dbi)
	if err != nil {
		return nil, err
	}

	//根据筛选方式，筛选出待备份的数据库
	var preDBS []string
	var errDBS []string = nil
	var allFlag = false
	for _, v := range *dbs {
		var errFlag = true
		for _, w := range *allDbs {
			v = strings.TrimSpace(v)
			w = strings.TrimSpace(w)
			if v == "all" {
				preDBS = nil
				preDBS = append(preDBS, "all")
				allFlag = true
				errFlag = false
				break
			} else if filterMethod(v, w, conf.FILTER_METHOD) {
				preDBS = append(preDBS, v)
				errFlag = false
				break
			}
		}

		if errFlag {
			errDBS = append(errDBS, v)
		}

		if allFlag {
			break
		}
	}

	if !allFlag {
		for _, v := range errDBS {
			logging.Logger.Printf("数据库名 %v 不存在", v)
		}
	}

	return &preDBS, nil
}

func filterMethod(w, v string, flag bool) bool {
	if flag {
		if w == v {
			return true
		} else {
			return false
		}
	} else {
		if w != v {
			return true
		} else {
			return false
		}
	}
}
