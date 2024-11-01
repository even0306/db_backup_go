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

	_ "github.com/go-sql-driver/mysql"
)

// 使用mysqldump备份mysql数据库，传入DBInfo结构体和要备份的数据库名指针，返回错误
func MysqlDump(info *DBInfo, db string, dst string, filename string, single int) error {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", info.DBVersion)
	if flag {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "--column-statistics=0", "-E", "-R", "--triggers", "--single-transaction", "--hex-blob", db)
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "--column-statistics=0", "-E", "-R", "--triggers", "--hex-blob", db)
		}
	} else {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "-E", "-R", "--triggers", "--single-transaction", "--hex-blob", db)
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-q", "-E", "-R", "--triggers", "--hex-blob", db)
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

	var gz *gzip.Writer
	reader := bufio.NewReader(stdout)
	isFirst := true
	//实时循环读取输出流中的一行内容
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
			logging.Logger.Panic(err)
			break
		}
		_, err = gz.Write([]byte(line)) //写入文件(字节数组)
		if err != nil {
			logging.Logger.Panic(err)
		}
	}
	return err

}

// 使用mysqldump备份mysql数据库，传入DBInfo结构体，返回错误
func MysqlDumpAll(info *DBInfo, dst string, filename string, single int) error {
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

	var gz *gzip.Writer
	reader := bufio.NewReader(stdout)
	isFirst := true
	//实时循环读取输出流中的一行内容
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
			logging.Logger.Panic(err)
			break
		}
		_, err = gz.Write([]byte(line)) //写入文件(字节数组)
		if err != nil {
			logging.Logger.Panic(err)
		}
	}
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
