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
	"regexp"
	"strings"
)

// 使用mysqldump备份mysql数据库，传入DBInfo结构体和要备份的数据库名指针，返回错误
func MysqlDump(info *DBInfo, db *string, dst string, filename string, single int) error {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", info.DBVersion)
	if flag {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "--column-statistics=0", "-E", "-R", "--triggers", "--single-transaction", *db)
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "--column-statistics=0", "-E", "-R", "--triggers", *db)
		}
	} else {
		if single == 1 {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-E", "-R", "--triggers", "--single-transaction", *db)
		} else {
			cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-E", "-R", "--triggers", *db)
		}
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		br := bufio.NewReader(strings.NewReader(stderr.String()))
		t, _, err := br.ReadLine()
		if err != nil {
			return err
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

	err = os.MkdirAll(dst+"/"+*db, 0777)
	if err != nil {
		return fmt.Errorf("创建备份文件路径失败：%w", err)
	}

	f, err := os.Create(dst + "/" + *db + "/" + filename)
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

// 使用mysqldump备份mysql数据库，传入DBInfo结构体，返回错误
func MysqlDumpAll(info *DBInfo, dst string, filename string, single int) error {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", info.DBVersion)
	if flag {
		cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "--column-statistics=0", "-E", "-R", "--triggers", "--skip-lock-tables", "--all-databases")
	} else {
		cmd = exec.Command(info.ExecPath+"/mysqldump", "-h"+info.DBHost, "-P"+fmt.Sprint(info.DBPort), "-u"+info.DBUser, "-p"+info.DBPassword, "-E", "-R", "--triggers", "--skip-lock-tables", "--all-databases")
	}

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
