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
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetMysqlDBList(t *testing.T) {
	var info struct {
		DBUser     string
		DBHost     string
		DBPassword string
		DBPort     int
	}
	info.DBUser = "root"
	info.DBHost = "127.0.0.1"
	info.DBPassword = "123456"
	info.DBPort = 3306

	db, err := sql.Open("mysql", info.DBUser+":"+info.DBPassword+"@tcp("+info.DBHost+":"+fmt.Sprint(info.DBPort)+")/information_schema?charset=utf8")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	rows, err := db.Query("select schema_name from schemata")
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()
	if err != nil {
		t.Error(err)
	}
	var list []string
	for rows.Next() {
		var col string
		rows.Scan(&col)
		list = append(list, col)
	}
	for _, v := range list {
		fmt.Println(v)
	}
}

func TestMysqlDumpAll(t *testing.T) {
	var info struct {
		DBUser     string
		DBHost     string
		DBPassword string
		DBPort     int
		DBVersion  string
		ExecPath   string
	}
	info.DBUser = "root"
	info.DBHost = "127.0.0.1"
	info.DBPassword = "1234561"
	info.DBPort = 3306
	info.DBVersion = "5.7"
	info.ExecPath = "/usr/bin"
	single := 1
	dst := "/app/mysql_backup"
	filename := "abc.sql.gz"

	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", info.DBVersion)
	if flag {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "--column-statistics=0", "-E", "-R", "--triggers", "--single-transaction", "--all-databases")
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "--column-statistics=0", "-E", "-R", "--triggers", "--all-databases")
		}
	} else {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "-E", "-R", "--triggers", "--single-transaction", "--all-databases")
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "-E", "-R", "--triggers", "--all-databases")
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
			t.Fatal(err)
		}
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		logging.Logger.Panic(err)
	}

	var gz *gzip.Writer
	reader := bufio.NewReader(stdout)
	isFirst := true
	//实时循环读取输出流中的一行内容
	for {
		line, err := reader.ReadString('\n')
		if isFirst {
			if line == "" {
				cmd.Stderr = os.Stderr
				t.Fatal(stderr.String())
			} else {
				err = os.MkdirAll(dst+"/all", 0777)
				if err != nil {
					t.Fatal(err)
				}

				f, err := os.Create(dst + "/all/" + filename)
				if err != nil {
					t.Fatal(err)
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
			t.Fatal(err)
			break
		}
		_, err = gz.Write([]byte(line)) //写入文件(字节数组)
		if err != nil {
			t.Fatal(err)
		}
	}
}
