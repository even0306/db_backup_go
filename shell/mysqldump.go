package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// 使用mysqldump备份mysql数据库，传入DBInfo结构体和要备份的数据库名指针，返回备份出的[]byte数据指针和错误
func MysqlDump(info *DBInfo, db *string) (*[]byte, error) {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", info.DBVersion)
	if flag && info.dbType == "mysql" {
		cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "--column-statistics=0", "-E", "-R", "--triggers", "--skip-lock-tables", *db)
	} else {
		cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-E", "-R", "--triggers", "--skip-lock-tables", *db)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		br := bufio.NewReader(strings.NewReader(stderr.String()))
		t, _, err := br.ReadLine()
		if err != nil {
			return nil, err
		}
		warn, err := regexp.Match("[Warning]", t)
		if warn {
			log.Print(string(t))
		} else {
			return nil, fmt.Errorf(*db+" 数据库备份失败：%w:%v", err, stderr.String())
		}
	}
	out := stdout.Bytes()
	return &out, nil
}

// 使用mysqldump备份mysql数据库，传入DBInfo结构体，返回备份出的[]byte数据指针和错误
func MysqlDumpAll(info *DBInfo) (*[]byte, error) {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", info.DBVersion)
	if flag && info.dbType == "mysql" {
		cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "--column-statistics=0", "-E", "-R", "--triggers", "--skip-lock-tables", "--all-databases")
	} else {
		cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-E", "-R", "--triggers", "--skip-lock-tables", "--all-databases")
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		br := bufio.NewReader(strings.NewReader(stderr.String()))
		t, _, err := br.ReadLine()
		if err != nil {
			return nil, err
		}
		warn, err := regexp.Match("[Warning]", t)
		if warn {
			log.Print(string(t))
		} else {
			return nil, fmt.Errorf("ALL数据库备份失败：%w:%v", err, stderr.String())
		}
	}
	out := stdout.Bytes()
	return &out, nil
}

// 使用mysql客户端查看mysql数据库现有的库，返回*[]string的数据库列表切片指针
func GetMysqlDBList(info *DBInfo) (*[]string, error) {
	cmd := exec.Command(info.ExecPath+"/mysql", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-Bse", "show databases")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("数据库列表查询失败：%w:%v", err, stderr.String())
	}
	out := stdout.String()
	list := strings.Split(out, "\n")
	return &list, nil
}
