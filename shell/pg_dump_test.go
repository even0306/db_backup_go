package shell

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"db_backup_go/logging"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	_ "github.com/lib/pq"
)

// 使用postgresql客户端查看postgresql数据库现有的库，返回*[]string的数据库列表切片指针
func TestPostgresqlDumpAll(t *testing.T) {
	var info struct {
		DBUser     string
		DBHost     string
		DBPassword string
		DBPort     int
		ExecPath   string
	}
	info.DBUser = "postgres"
	info.DBHost = "127.0.0.1"
	info.DBPassword = "123456"
	info.DBPort = 5432
	info.ExecPath = "/usr/bin"
	dst := "/app/mysql_backup/"
	filename := "abc.sql.gz"

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
				t.Fatal(stderr.String())
			} else {
				err = os.MkdirAll(dst+"/all", 0777)
				if err != nil {
					t.Error(err)
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
}
