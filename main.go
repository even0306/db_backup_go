package main

import (
	"db_backup_go/controller"
	"flag"
	"log"
	"os"
	"path/filepath"
)

//程序主入口
func main() {
	version := "1.1.8"
	printVersion := flag.Bool("version", false, "[--version]")
	flag.Parse()
	if *printVersion {
		println(version)
		os.Exit(0)
	}

	logfile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	log.SetOutput(logfile)
	if err != nil {
		log.Panicf("打开日志文件失败：%v", err)
	}

	ex, err := os.Executable()
	exPath := filepath.Dir(ex)
	_, err = os.Stat(exPath + "/config.json")
	if err == nil {
		_, err = os.Stat(exPath + "/dbs.txt")
		if err == nil {
			c := controller.NewController(exPath+"/config.json", exPath+"/dbs.txt")
			err = c.Controller()
			if err != nil {
				log.Panic(err)
			}
		}
		if os.IsNotExist(err) {
			log.Panic("找不到dbs.txt")
		}
	}
	if os.IsNotExist(err) {
		log.Panic("找不到config.json")
	}
}
