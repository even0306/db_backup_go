package shell

import "db_backup_go/logging"

type DBInfo struct {
	dbType            string
	ExecPath          string
	DBVersion         string
	DBHost            string
	DBPort            int
	DBUser            string
	DBPassword        string
	singleTransaction int
}

func NewSelecter(dbType string, p string, ver string, host string, port int, user string, pass string, single int) *DBInfo {
	return &DBInfo{
		dbType:            dbType,
		ExecPath:          p,
		DBVersion:         ver,
		DBHost:            host,
		DBPort:            port,
		DBUser:            user,
		DBPassword:        pass,
		singleTransaction: single,
	}
}

// 备份工具选择器，传入 *common.ConfigFile 和要备份的库名指针，返回备份出的字节流指针和报错信息
func BackupSelecter(b *DBInfo, db *string, dst string, filename *string) error {
	var err error
	if b.dbType == "mysql" || b.dbType == "mariadb" {
		if *db == "all" {
			err = MysqlDumpAll(b, dst, *filename)
			if err != nil {
				return err
			}
		} else {
			err = MysqlDump(b, db, dst, *filename)
			if err != nil {
				return err
			}
		}
	} else if b.dbType == "postgresql" {
		if *db == "all" {
			err = PostgresqlDumpAll(b, dst, *filename)
			if err != nil {
				return err
			}
		} else {
			err = PostgresqlDump(b, db, dst, *filename)
			if err != nil {
				return err
			}
		}
	} else {
		logging.Logger.Panic("未知的数据库类型，请重新检查 config.json 文件配置")
	}
	return nil
}

func DBListSelecter(b *DBInfo) (*[]string, error) {
	if b.dbType == "mysql" || b.dbType == "mariadb" {
		allDbs, err := GetMysqlDBList(b)
		if err != nil {
			return nil, err
		}
		return allDbs, nil
	} else if b.dbType == "postgresql" {
		allDbs, err := GetPostgresqlDBList(b)
		if err != nil {
			return nil, err
		}
		return allDbs, nil
	} else {
		logging.Logger.Panic("未知的数据库类型，请重新检查 config.json 文件配置")
	}
	return nil, nil
}
