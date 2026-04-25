package shell

import (
	"bufio"
	"bytes"
	"database/sql"
	"db_backup_go/config"
	"db_backup_go/logging"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 使用mysqldump备份mysql数据库，传入DBInfo结构体和要备份的数据库名指针，返回错误
func MysqlDump(info *DBInfo, savePath config.BackupPath, single int, dbBackupReference string) error {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", info.DBVersion)
	if flag {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "--column-statistics=0", "-E", "-R", "--triggers", "--single-transaction", "--hex-blob", dbBackupReference)
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "--column-statistics=0", "-E", "-R", "--triggers", "--hex-blob", dbBackupReference)
		}
	} else {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "-E", "-R", "--triggers", "--single-transaction", "--hex-blob", dbBackupReference)
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "-E", "-R", "--triggers", "--hex-blob", dbBackupReference)
		}
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		br := bufio.NewReader(strings.NewReader(stderr.String()))
		a, _, err := br.ReadLine()
		if err != nil {
			logging.Logger.Panic(err)
		}
		warn, err := regexp.Match("Warning", a)
		if warn {
			logging.Logger.Print(string(a))
		} else {
			logging.Logger.Panic(err)
		}
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		logging.Logger.Panic(err)
	}

	WriteToFile(&stdout, &stderr, savePath)

	return err

}

// 使用mysqldump备份mysql数据库，传入DBInfo结构体，返回错误
func MysqlDumpAll(info *DBInfo, savePath config.BackupPath, single int) error {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", info.DBVersion)
	if flag {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "--column-statistics=0", "-E", "-R", "--triggers", "--single-transaction", "--hex-blob", "--all-databases")
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "--column-statistics=0", "-E", "-R", "--triggers", "--hex-blob", "--all-databases")
		}
	} else {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "-E", "-R", "--triggers", "--single-transaction", "--hex-blob", "--all-databases")
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "-E", "-R", "--triggers", "--hex-blob", "--all-databases")
		}
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		br := bufio.NewReader(strings.NewReader(stderr.String()))
		a, _, err := br.ReadLine()
		if err != nil {
			logging.Logger.Panic(err)
		}
		warn, err := regexp.Match("Warning", a)
		if warn {
			logging.Logger.Print(string(a))
		} else {
			logging.Logger.Panic(err)
		}
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		logging.Logger.Panic(err)
	}

	WriteToFile(&stdout, &stderr, savePath)

	return err
}

// 使用mysql客户端查看mysql数据库现有的库，返回*[]string的数据库列表切片指针
func GetMysqlDBList(info *DBInfo) (*[]string, error) {
	db, err := sql.Open("mysql", info.DBUser+":"+info.DBPassword+"@tcp("+info.DBHost+":"+fmt.Sprint(info.DBPort)+")/information_schema?charset=utf8")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select schema_name from schemata")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
