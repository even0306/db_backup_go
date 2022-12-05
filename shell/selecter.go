package shell

type DBInfo struct {
	dbType     string
	ExecPath   string
	DBVersion  string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
}

func NewSelecter(dbType string, p string, ver string, host string, port int, user string, pass string) *DBInfo {
	return &DBInfo{
		dbType:     dbType,
		ExecPath:   p,
		DBVersion:  ver,
		DBHost:     host,
		DBPort:     port,
		DBUser:     user,
		DBPassword: pass,
	}
}

//备份工具选择器，传入 *common.ConfigFile 和要备份的库名指针，返回备份出的字节流指针和报错信息
func BackupSelecter(b *DBInfo, db *string) (*[]byte, error) {
	var out *[]byte
	var err error
	if b.dbType == "mysql" || b.dbType == "mariadb" {
		if *db == "all" {
			out, err = MysqlDumpAll(b)
			if err != nil {
				return nil, err
			}
		} else {
			out, err = MysqlDump(b, db)
			if err != nil {
				return nil, err
			}
		}
	} else if b.dbType == "postgresql" {
		if *db == "all" {
			out, err = PostgresqlDumpAll(b)
			if err != nil {
				return nil, err
			}
		} else {
			out, err = PostgresqlDump(b, db)
			if err != nil {
				return nil, err
			}
		}
	}
	return out, nil
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
	}
	return nil, nil
}
