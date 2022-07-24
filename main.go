package main

import (
	"db_backup_go/controller"
	"flag"
	"log"
	"os"
)

//程序主入口
func main() {
	version := "1.1.4"
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

	_, err = os.Stat("config.json")
	if err == nil {
		_, err = os.Stat("dbs.txt")
		if err == nil {
			c := controller.NewController("config.json", "dbs.txt")
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
