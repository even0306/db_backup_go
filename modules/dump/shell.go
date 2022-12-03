package dump

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type DBDump interface {
	MysqlDump(db *string) ([]byte, error)
	MysqlDumpAll() ([]byte, error)
	GetMysqlDBList() (*[]string, error)
	PostgresqlDump(db *string) ([]byte, error)
	PostgresqlDumpAll() ([]byte, error)
	GetPostgresqlDBList() (*[]string, error)
}

type dbDump struct {
	DBInfo
	dumpExecPath string
}

type DBInfo struct {
	DBVersion  string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
}

//初始化备份方法，传入备份命令所在路径，*DBInfo类型指针，里面包含数据库连接信息
func NewDBDumpFunc(dumpExecPath string, dbi *DBInfo) *dbDump {
	return &dbDump{
		DBInfo: DBInfo{
			DBVersion:  "",
			DBHost:     dbi.DBHost,
			DBPort:     dbi.DBPort,
			DBUser:     dbi.DBUser,
			DBPassword: dbi.DBPassword,
		},
		dumpExecPath: dumpExecPath,
	}
}

var stdout bytes.Buffer
var stderr bytes.Buffer

// 使用mysqldump备份mysql数据库，传入DBInfo结构体和要备份的数据库名指针，返回备份出的[]byte数据指针和错误
func (d *dbDump) MysqlDump(db *string) (*[]byte, error) {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", d.DBVersion)
	if flag {
		cmd = exec.Command(d.dumpExecPath+"/mysqldump", "-h"+d.DBHost, "-P"+fmt.Sprint(d.DBPort), "-u"+d.DBUser, "-p"+d.DBPassword, "--column-statistics=0", "-E", "-R", "--triggers", "--skip-lock-tables", *db)
	} else {
		cmd = exec.Command(d.dumpExecPath+"/mysqldump", "-h"+d.DBHost, "-P"+fmt.Sprint(d.DBPort), "-u"+d.DBUser, "-p"+d.DBPassword, "-E", "-R", "--triggers", "--skip-lock-tables", *db)
	}
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
func (d *dbDump) MysqlDumpAll() (*[]byte, error) {
	var cmd *exec.Cmd
	flag, _ := regexp.MatchString("8.0.*", d.DBVersion)
	if flag {
		cmd = exec.Command(d.dumpExecPath+"/mysqldump", "-h"+d.DBHost, "-P"+fmt.Sprint(d.DBPort), "-u"+d.DBUser, "-p"+d.DBPassword, "--column-statistics=0", "-E", "-R", "--triggers", "--skip-lock-tables", "--all-databases")
	} else {
		cmd = exec.Command(d.dumpExecPath+"/mysqldump", "-h"+d.DBHost, "-P"+fmt.Sprint(d.DBPort), "-u"+d.DBUser, "-p"+d.DBPassword, "-E", "-R", "--triggers", "--skip-lock-tables", "--all-databases")
	}
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
func (d *dbDump) GetMysqlDBList() (*[]string, error) {
	cmd := exec.Command(d.dumpExecPath+"/mysql", "-h"+d.DBHost, "-P"+fmt.Sprint(d.DBPort), "-u"+d.DBUser, "-p"+d.DBPassword, "-Bse", "show databases")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ALL数据库备份失败：%w:%v", err, stderr.String())
	}
	out := stdout.String()
	list := strings.Split(out, "\n")
	return &list, nil
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体和要备份的数据库名指针，返回备份出的[]byte数据指针和错误
func (d *dbDump) PostgresqlDump(db *string) (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/pg_dump", "-h", d.DBHost, "-p", fmt.Sprint(d.DBPort), "-U", d.DBUser, "-d", *db, "--inserts")

	env := os.Environ()
	cmdEnv := []string{}
	flag := false
	for _, e := range env {
		i := strings.Index(e, "=")
		if i > 0 && (e[:i] == "PGPASSWORD") {
			e = "PGPASSWORD=" + d.DBPassword
			cmdEnv = append(cmdEnv, e)
			flag = true
			break
		} else {
			cmdEnv = append(cmdEnv, e)
		}
	}
	if !flag {
		cmdEnv = append(cmdEnv, "PGPASSWORD="+d.DBPassword)
	}
	cmd.Env = cmdEnv
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
func (d *dbDump) PostgresqlDumpAll() (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/pg_dumpall", "-h", d.DBHost, "-p", fmt.Sprint(d.DBPort), "-U", d.DBUser, "--inserts")

	env := os.Environ()
	cmdEnv := []string{}
	flag := false
	for _, e := range env {
		i := strings.Index(e, "=")
		if i > 0 && (e[:i] == "PGPASSWORD") {
			e = "PGPASSWORD=" + d.DBPassword
			cmdEnv = append(cmdEnv, e)
			flag = true
			break
		} else {
			cmdEnv = append(cmdEnv, e)
		}
	}
	if !flag {
		cmdEnv = append(cmdEnv, "PGPASSWORD="+d.DBPassword)
	}
	cmd.Env = cmdEnv
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
func (d *dbDump) GetPostgresqlDBList() (*[]string, error) {
	cmd := exec.Command(d.dumpExecPath+"/psql", fmt.Sprintf("host=%s port=%v user=%s password=%s", d.DBHost, d.DBPort, d.DBUser, d.DBPassword), "-c", "SELECT datname FROM pg_database;")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("stderr: %w", err)
	}
	list := strings.Split(string(out), "\n")
	list = list[2 : len(list)-3]
	return &list, nil
}
