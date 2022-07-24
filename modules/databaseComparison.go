package modules

import (
	"db_backup_go/common"
	"log"
	"strings"
)

type Comparison interface {
	Comparison() (*[]string, error)
}

type comparisonInfo struct {
	conf *common.ConfigFile
	dbs  *[]string
}

func NewCompartor(conf *common.ConfigFile, dbs *[]string) *comparisonInfo {
	return &comparisonInfo{
		conf: conf,
		dbs:  dbs,
	}
}

func (c *comparisonInfo) Comparison() (*[]string, error) {
	//获取所有数据库名
	var err error
	var allDbs *[]string
	dbi := DBInfo{
		DBVersion:  c.conf.DB_Version,
		DBHost:     c.conf.DB_HOST,
		DBPort:     c.conf.DB_PORT,
		DBUser:     c.conf.DB_USER,
		DBPassword: c.conf.DB_PASSWORD,
	}
	dbu := NewDBDumpFunc(c.conf.MYSQL_EXEC_PATH, &dbi)
	if c.conf.DATABASETYPE == "mysql" {
		allDbs, err = dbu.GetMysqlDBList()
		if err != nil {
			return nil, err
		}
	} else if c.conf.DATABASETYPE == "postgresql" {
		allDbs, err = dbu.GetPostgresqlDBList()
		if err != nil {
			return nil, err
		}
	}

	//根据筛选方式，筛选出待备份的数据库
	var preDBS []string
	var errDBS []string = nil
	var allFlag = false
	for _, v := range *c.dbs {
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
			} else if filterMethod(v, w, c.conf.FILTER_METHOD) {
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

	for _, v := range errDBS {
		log.Printf("数据库 %v 不存在，备份失败", v)
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
