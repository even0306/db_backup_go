package shell

import (
	"bytes"
	"database/sql"
	"db_backup_go/config"
	"db_backup_go/logging"
	"fmt"
	"os/exec"

	_ "github.com/lib/pq"
)

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体和要备份的数据库名指针，返回错误
func PostgresqlDump(dbInfo *DBInfo, savePath config.BackupPath, db string) error {
	cmd := exec.Command(dbInfo.ExecPath+"/pg_dump", "-h", dbInfo.DBHost, "-p", fmt.Sprint(dbInfo.DBPort), "-U", dbInfo.DBUser, "-d", db, "--inserts")
	cmd.Env = []string{"PGPASSWORD=" + dbInfo.DBPassword}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logging.Logger.Panic(err)
	}

	err = cmd.Start()
	if err != nil {
		logging.Logger.Panic(err)
	}

	WriteToFile(&stdout, &stderr, savePath)

	return err
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体，返回错误
func PostgresqlDumpAll(dbInfo *DBInfo, savePath config.BackupPath) error {
	cmd := exec.Command(dbInfo.ExecPath+"/pg_dumpall", "-h", dbInfo.DBHost, "-p", fmt.Sprint(dbInfo.DBPort), "-U", dbInfo.DBUser, "--inserts")
	cmd.Env = []string{"PGPASSWORD=" + dbInfo.DBPassword}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logging.Logger.Panic(err)
	}

	err = cmd.Start()
	if err != nil {
		logging.Logger.Panic(err)
	}

	WriteToFile(&stdout, &stderr, savePath)

	return err
}

// 使用postgresql客户端查看postgresql数据库现有的库，返回*[]string的数据库列表切片指针
func GetPostgresqlDBList(dbInfo *DBInfo) (*[]string, error) {
	db, err := sql.Open("postgres", "host="+dbInfo.DBHost+" port="+fmt.Sprint(dbInfo.DBPort)+" user="+dbInfo.DBUser+" password="+dbInfo.DBPassword+" dbname=postgres"+" sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select datname from pg_catalog.pg_database")
	if err != nil {
		return nil, err
	}

	var list []string
	for rows.Next() {
		var col string
		err = rows.Scan(&col)
		if err != nil {
			return nil, err
		}
		list = append(list, col)
	}
	return &list, nil
}
