package main

import (
	"db_backup_go/controller"
	"db_backup_go/logging"
	"embed"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
)

//go:embed assets
var f embed.FS

// 程序主入口
func main() {
	version := "1.2.0"
	printVersion := flag.Bool("version", false, "[--version]")
	flag.Parse()
	if *printVersion {
		println(version)
		os.Exit(0)
	}

	//找到执行程序所在位置
	ex, err := os.Executable()
	if err != nil {
		logging.Logger.Panic(err)
	}
	exPath := filepath.Dir(ex)

	//初始化日志
	logging.NewLogger(exPath + "/server.log")

	_, err = os.Stat(exPath + "/config.json")
	if err != nil {
		fv, err := f.ReadFile("assets/config.json")
		if err != nil {
			panic(err)
		}
		rf, err := os.Create("config.json")
		if err != nil {
			panic(err)
		}
		defer rf.Close()
		io.WriteString(rf, string(fv))
		log.Panic("未找到配置文件 config.json，已创建默认配置，请进行修改后重新执行。")
	}
	_, err = os.Stat(exPath + "/dbs.txt")
	if err != nil {
		fv, err := f.ReadFile("assets/dbs.txt")
		if err != nil {
			panic(err)
		}
		rf, err := os.Create("dbs.txt")
		if err != nil {
			panic(err)
		}
		defer rf.Close()
		io.WriteString(rf, string(fv))
		log.Panic("未找到配置文件 dbs.txt，已创建默认配置，默认备份所有库。")
	}

	c := controller.NewController(exPath+"/config.json", exPath+"/dbs.txt")
	err = c.Controller()
	if err != nil {
		logging.Logger.Panic(err)
	}
}
