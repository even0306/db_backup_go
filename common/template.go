package common

import (
	"os"
)

func CreateConfig() error {
	defaultConfig := `{
	# 目前支持 mysql 和 postgresql
	"database": "postgresql",

	# true 正向筛选，false 反向筛选
	"filter_method": true,
	"SAVE_DAY": 7,
	"mysql_exec_path": "/usr/bin/",
	"backup_save_path": "/app/mysql_backup/",

	"DB_VERSION": "14.0",
	"DB_HOST": "127.0.0.1",
	"DB_PORT": 5432,
	"DB_USER": "postgres",
	"DB_PASSWORD": "123456",
	"DB_LABEL": "db_16",

	# true 开启发送到异机功能，false 关闭发送到异机功能
	"remote_backup": true,
	"REMOTE_HOST": "192.168.56.15",
	"REMOTE_PORT": 22,
	"REMOTE_USER": "even",
	"REMOTE_PASSWORD": "2333",
	"REMOTE_PATH": "/app/backup"
}
`

	dc, err := os.OpenFile("./config.json", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer dc.Close()

	_, err = dc.WriteString(defaultConfig)
	if err != nil {
		return err
	}

	return nil
}

func CreateDBS() error {
	defaultDbs := "all"

	dd, err := os.OpenFile("./dbs.txt", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer dd.Close()

	_, err = dd.WriteString(defaultDbs)
	if err != nil {
		return err
	}

	return nil
}
