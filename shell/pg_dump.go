package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体和要备份的数据库名指针，返回备份出的[]byte数据指针和错误
func PostgresqlDump(info *DBInfo, db *string) (*[]byte, error) {
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf(*db+" 数据库备份失败：%w:%v", err, stderr.String())
	}
	out := stdout.Bytes()
	return &out, nil
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体，返回备份出的[]byte数据指针和错误
func PostgresqlDumpAll(info *DBInfo) (*[]byte, error) {
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ALL数据库备份失败：%w:%v", err, stderr.String())
	}
	out := stdout.Bytes()
	return &out, nil
}

// 使用postgresql客户端查看postgresql数据库现有的库，返回*[]string的数据库列表切片指针
func GetPostgresqlDBList(info *DBInfo) (*[]string, error) {
	cmd := exec.Command(info.ExecPath+"/psql", fmt.Sprintf("host=%s port=%v user=%s password=%s", info.DBHost, info.DBPort, info.DBUser, info.DBPassword), "-c", "SELECT datname FROM pg_database;")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("数据库列表查询失败：%w:%v", err, stderr.String())
	}
	out := stdout.String()
	list := strings.Split(string(out), "\n")
	list = list[2 : len(list)-3]
	return &list, nil
}
