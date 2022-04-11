package modules

import (
	"fmt"
	"os/exec"
	"strconv"

	_ "github.com/bmizerany/pq"
)

type DBDump interface {
	MysqlDump(db *string) ([]byte, error)
	MysqlDumpAll() ([]byte, error)
	GetMysqlDBList() (*[]byte, error)
	PostgresqlDump(db *string) ([]byte, error)
	PostgresqlDumpAll() ([]byte, error)
	GetPostgresqlDBList() (*[]byte, error)
}

type dbDump struct {
	DBInfo
	dumpExecPath string
}

type DBInfo struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
}

func NewDBDumpFunc(dumpExecPath string, dbi *DBInfo) *dbDump {
	return &dbDump{
		DBInfo: DBInfo{
			DBHost:     dbi.DBHost,
			DBPort:     dbi.DBPort,
			DBUser:     dbi.DBUser,
			DBPassword: dbi.DBPassword,
		},
		dumpExecPath: dumpExecPath,
	}
}

// 使用mysqldump备份mysql数据库，传入DBInfo结构体和要备份的数据库名指针，返回备份出的[]byte数据指针和错误
func (d *dbDump) MysqlDump(db *string) (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/mysqldump", "-h"+d.DBHost, "-P"+strconv.Itoa(d.DBPort), "-u"+d.DBUser, "-p"+d.DBPassword, "-E", "-R", "--triggers", *db)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(*db+"数据库配置失败：%w", err)
	}
	return &out, nil
}

// 使用mysqldump备份mysql数据库，传入DBInfo结构体，返回备份出的[]byte数据指针和错误
func (d *dbDump) MysqlDumpAll() (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/mysqldump", "-h"+d.DBHost, "-P"+strconv.Itoa(d.DBPort), "-u"+d.DBUser, "-p"+d.DBPassword, "-E", "-R", "--triggers", "--all-databases")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("all数据库配置失败：%w", err)
	}
	return &out, nil
}

// 使用mysql客户端查看mysql数据库现有的库，返回[]byte的数据库列表
func (d *dbDump) GetMysqlDBList() (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/mysql", "-h"+d.DBHost, "-P"+strconv.Itoa(d.DBPort), "-u"+d.DBUser, "-p"+d.DBPassword, "-Bse", "show databases")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("stderr: %w", err)
	}
	return &out, nil
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体和要备份的数据库名指针，返回备份出的[]byte数据指针和错误
func (d *dbDump) PostgresqlDump(db *string) (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/pg_dump", "\"host="+d.DBHost, "port="+strconv.Itoa(d.DBPort), "user="+d.DBUser, "password="+d.DBPassword+"\"", "-t", *db, "--inserts")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(*db+"数据库配置失败：%w", err)
	}
	return &out, nil
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体，返回备份出的[]byte数据指针和错误
func (d *dbDump) PostgresqlDumpAll() (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/pg_dumpall", "\"host="+d.DBHost, "port="+strconv.Itoa(d.DBPort), "user="+d.DBUser, "password="+d.DBPassword+"\"", "--inserts")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("all数据库配置失败：%w", err)
	}
	return &out, nil
}

func (d *dbDump) GetPostgresqlDBList() (*[]byte, error) {
	cmd1 := exec.Command(d.dumpExecPath+"/psql", "`host="+d.DBHost+" port="+strconv.Itoa(d.DBPort)+" user="+d.DBUser+" password="+d.DBPassword+"`", "-c", "\"SELECT datname FROM pg_database;\"")
	// cmd2 := exec.Command("awk", "'{if (NR>2){print $1}}'")
	// cmd3 := exec.Command("awk", "'NR>1 {print last} {last=$0}'")
	// cmd4 := exec.Command("awk", "'NR>1 {print last} {last=$0}'")
	// cmd2.Stdin, _ = cmd1.StdoutPipe()
	// cmd3.Stdin, _ = cmd2.StdoutPipe()
	// cmd4.Stdin, _ = cmd3.StdoutPipe()
	out, err := cmd1.Output()
	if err != nil {
		return nil, fmt.Errorf("stderr: %w", err)
	}
	return &out, nil
}
