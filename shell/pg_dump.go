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
	"regexp"
	"strings"

	_ "github.com/lib/pq"
)

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体和要备份的数据库名指针，返回错误
func PostgresqlDump(info *DBInfo, db *string, dst string, filename string) error {
	cmd := exec.Command(info.ExecPath+"/pg_dump", "-h", info.DBHost, "-p", fmt.Sprint(info.DBPort), "-U", info.DBUser, "-d", *db, "--inserts")
	env := os.Environ()
	cmdEnv := []string{}
	flag := false
	for _, e := range env {
		i := strings.Index(e, "=")
		if i > 0 && (e[:i] == "PGPASSWORD") {
			e = "PGPASSWORD=" + info.DBPassword
			cmdEnv = append(cmdEnv, e)
			flag = true
			break
		} else {
			cmdEnv = append(cmdEnv, e)
		}
	}
	if !flag {
		cmdEnv = append(cmdEnv, "PGPASSWORD="+info.DBPassword)
	}
	cmd.Env = cmdEnv

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		br := bufio.NewReader(strings.NewReader(stderr.String()))
		t, _, err := br.ReadLine()
		if err != nil {
			logging.Logger.Panic(err)
		}
		warn, err := regexp.Match("[Warning]", t)
		if warn {
			logging.Logger.Print(string(t))
		} else {
			return fmt.Errorf(*db+" 数据库备份失败：%w:%v", err, stderr.String())
		}
	}
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		logging.Logger.Panic(err)
	}

	err = os.MkdirAll(dst+"/all", 0777)
	if err != nil {
		return fmt.Errorf("创建备份文件路径失败：%w", err)
	}

	f, err := os.Create(dst + "/all/" + filename)
	if err != nil {
		logging.Logger.Panic(err)
	}
	defer f.Close()

	//创建一个gzip的流来接收管道中内容
	gz := gzip.NewWriter(f)
	defer gz.Close()

	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			logging.Logger.Panicf("读取流出现问题：%v，文件备份不完整。", err)
			break
		}
		_, err = gz.Write([]byte(line)) //写入文件(字节数组)
		if err != nil {
			logging.Logger.Panic(err)
		}
		f.Sync()
	}
	err = cmd.Wait()
	return err
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体，返回错误
func PostgresqlDumpAll(info *DBInfo, dst string, filename string) error {
	cmd := exec.Command(info.ExecPath+"/pg_dumpall", "-h", info.DBHost, "-p", fmt.Sprint(info.DBPort), "-U", info.DBUser, "--inserts")
	env := os.Environ()
	cmdEnv := []string{}
	flag := false
	for _, e := range env {
		i := strings.Index(e, "=")
		if i > 0 && (e[:i] == "PGPASSWORD") {
			e = "PGPASSWORD=" + info.DBPassword
			cmdEnv = append(cmdEnv, e)
			flag = true
			break
		} else {
			cmdEnv = append(cmdEnv, e)
		}
	}
	if !flag {
		cmdEnv = append(cmdEnv, "PGPASSWORD="+info.DBPassword)
	}
	cmd.Env = cmdEnv

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		br := bufio.NewReader(strings.NewReader(stderr.String()))
		t, _, err := br.ReadLine()
		if err != nil {
			logging.Logger.Panic(err)
		}
		warn, err := regexp.Match("[Warning]", t)
		if warn {
			logging.Logger.Print(string(t))
		} else {
			return fmt.Errorf("all 数据库备份失败：%w:%v", err, stderr.String())
		}
	}
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		logging.Logger.Panic(err)
	}

	err = os.MkdirAll(dst+"/all", 0777)
	if err != nil {
		return fmt.Errorf("创建备份文件路径失败：%w", err)
	}

	f, err := os.Create(dst + "/all/" + filename)
	if err != nil {
		logging.Logger.Panic(err)
	}
	defer f.Close()

	//创建一个gzip的流来接收管道中内容
	gz := gzip.NewWriter(f)
	defer gz.Close()

	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			logging.Logger.Panicf("读取流出现问题：%v，文件备份不完整。", err)
			break
		}
		_, err = gz.Write([]byte(line)) //写入文件(字节数组)
		if err != nil {
			logging.Logger.Panic(err)
		}
		f.Sync()
	}
	err = cmd.Wait()
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
