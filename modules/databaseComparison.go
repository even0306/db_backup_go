package modules

import (
	"bufio"
	"bytes"
	"db_backup_go/common"
	"errors"
	"fmt"
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
	var out *[]byte
	var err error
	dbi := DBInfo{
		DBHost:     c.conf.DB_HOST,
		DBPort:     c.conf.DB_PORT,
		DBUser:     c.conf.DB_USER,
		DBPassword: c.conf.DB_PASSWORD,
	}
	dbu := NewDBDumpFunc(c.conf.MYSQL_EXEC_PATH, &dbi)
	if c.conf.DATABASETYPE == "mysql" {
		out, err = dbu.GetMysqlDBList()
		if err != nil {
			return nil, err
		}
	} else if c.conf.DATABASETYPE == "postgresql" {
		out, err = dbu.GetPostgresqlDBList()
		if err != nil {
			return nil, err
		}
	}

	var allDbs []string

	rf := bytes.NewReader(*out)
	bf := bufio.NewScanner(rf)

	for {
		flag := bf.Scan()
		line := bf.Text()
		if err != nil && errors.Is(err, bf.Err()) == true {
			return nil, fmt.Errorf("逐行读取系统数据库名失败：%w", err)
		}
		if flag == false {
			break
		}
		allDbs = append(allDbs, line)
	}

	//根据筛选方式，筛选出待备份的数据库
	var preDBS []string
	var flag = false
	if c.conf.FILTER_METHOD == true {
		for _, v := range allDbs {
			for _, w := range *c.dbs {
				if w == "all" {
					preDBS = nil
					preDBS = append(preDBS, "all")
					flag = true
					break
				} else if w == string(v) {
					preDBS = append(preDBS, w)
				}
			}
			if flag == true {
				break
			}
		}
	} else {
		for _, v := range allDbs {
			for _, w := range *c.dbs {
				if w == "all" {
					preDBS = nil
					preDBS = append(preDBS, "all")
					flag = true
					break
				} else if w != string(v) {
					preDBS = append(preDBS, w)
				}
			}
			if flag == true {
				break
			}
		}
	}
	return &preDBS, nil
}
