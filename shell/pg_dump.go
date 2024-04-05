package shell

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"db_backup_go/logging"
	"fmt"
	"io"
	"os"
	"os/exec"

	_ "github.com/lib/pq"
)

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体和要备份的数据库名指针，返回错误
func PostgresqlDump(info *DBInfo, db string, dst string, filename string) error {
	cmd := exec.Command(info.ExecPath+"/pg_dump", "-h", info.DBHost, "-p", fmt.Sprint(info.DBPort), "-U", info.DBUser, "-d", db, "--inserts")
	cmd.Env = []string{"PGPASSWORD=" + info.DBPassword}

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

	var gz *gzip.Writer
	reader := bufio.NewReader(stdout)
	isFirst := true
	for {
		line, err := reader.ReadString('\n')
		if isFirst {
			if line == "" {
				cmd.Stderr = os.Stderr
				logging.Logger.Panic(stderr.String())
			} else {
				err = os.MkdirAll(dst+"/"+db, 0777)
				if err != nil {
					logging.Logger.Panic(err)
				}

				f, err := os.Create(dst + "/" + db + "/" + filename)
				if err != nil {
					logging.Logger.Panic(err)
				}
				defer f.Close()

				//创建一个gzip的流来接收管道中内容
				gz = gzip.NewWriter(f)
				defer gz.Close()
			}
			isFirst = false
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			logging.Logger.Panicf("读取流出现问题：%v，文件备份不完整。", err)
		}
		_, err = gz.Write([]byte(line))
		if err != nil {
			logging.Logger.Panic(err)
		}
	}
	return err
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体，返回错误
func PostgresqlDumpAll(info *DBInfo, dst string, filename string) error {
	cmd := exec.Command(info.ExecPath+"/pg_dumpall", "-h", info.DBHost, "-p", fmt.Sprint(info.DBPort), "-U", info.DBUser, "--inserts")
	cmd.Env = []string{"PGPASSWORD=" + info.DBPassword}

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

	var gz *gzip.Writer
	reader := bufio.NewReader(stdout)
	isFirst := true
	for {
		line, err := reader.ReadString('\n')
		if isFirst {
			if line == "" {
				cmd.Stderr = os.Stderr
				logging.Logger.Panic(stderr.String())
			} else {
				err = os.MkdirAll(dst+"/all", 0777)
				if err != nil {
					logging.Logger.Panic(err)
				}

				f, err := os.Create(dst + "/all/" + filename)
				if err != nil {
					logging.Logger.Panic(err)
				}
				defer f.Close()

				//创建一个gzip的流来接收管道中内容
				gz = gzip.NewWriter(f)
				defer gz.Close()
			}
			isFirst = false
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			logging.Logger.Panicf("读取流出现问题：%v，文件备份不完整。", err)
		}
		_, err = gz.Write([]byte(line))
		if err != nil {
			logging.Logger.Panic(err)
		}
	}
	return err
}

// 使用postgresql客户端查看postgresql数据库现有的库，返回*[]string的数据库列表切片指针
func GetPostgresqlDBList(info *DBInfo) (*[]string, error) {
	db, err := sql.Open("postgres", "host="+info.DBHost+" port="+fmt.Sprint(info.DBPort)+" user="+info.DBUser+" password="+info.DBPassword+" dbname=postgres"+" sslmode=disable")
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
