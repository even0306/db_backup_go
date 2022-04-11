package modules

import (
	"fmt"
	"os/exec"
)

type DBDump interface {
	MysqlDump(r DBInfo) ([]byte, error)
	MysqlDumpAll(r DBInfo) ([]byte, error)
	PostgresqlDump(r DBInfo) ([]byte, error)
	PostgresqlDumpAll(r DBInfo) ([]byte, error)
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

func NewDBDumpFunc(dumpExecPath string) *dbDump {
	return &dbDump{
		DBInfo:       DBInfo{},
		dumpExecPath: dumpExecPath,
	}
}

// 使用mysqldump备份mysql数据库，传入DBInfo结构体和要备份的数据库名指针，返回备份出的[]byte数据指针和错误
func (d *dbDump) MysqlDump(r DBInfo, db *string) (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/mysqldump", "-h"+d.DBHost, "-P"+string(rune(d.DBPort)), "-u"+d.DBUser, "-p"+d.DBPassword, "-E", "-R", "--triggers", *db)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(*db+"数据库配置失败：%w", err)
	}
	return &out, nil
}

// 使用mysqldump备份mysql数据库，传入DBInfo结构体，返回备份出的[]byte数据指针和错误
func (d *dbDump) MysqlDumpAll(r DBInfo) (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/mysqldump", "-h"+d.DBHost, "-P"+string(rune(d.DBPort)), "-u"+d.DBUser, "-p"+d.DBPassword, "-E", "-R", "--triggers", "--all-databases")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("all数据库配置失败：%w", err)
	}
	return &out, nil
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体和要备份的数据库名指针，返回备份出的[]byte数据指针和错误
func (d *dbDump) PostgresqlDump(r DBInfo, db *string) (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/pg_dump", "\"host="+d.DBHost, "port="+string(rune(d.DBPort)), "user="+d.DBUser, "password="+d.DBPassword+"\"", "-t", *db, "--inserts")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(*db+"数据库配置失败：%w", err)
	}
	return &out, nil
}

// 使用pg_dump备份postgresql数据库，传入DBInfo结构体，返回备份出的[]byte数据指针和错误
func (d *dbDump) PostgresqlDumpAll(r DBInfo) (*[]byte, error) {
	cmd := exec.Command(d.dumpExecPath+"/pg_dumpall", "\"host="+d.DBHost, "port="+string(rune(d.DBPort)), "user="+d.DBUser, "password="+d.DBPassword+"\"", "--inserts")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("all数据库配置失败：%w", err)
	}
	return &out, nil
}
