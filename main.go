package main

import (
	"flag"
	"log"
	"mysql_backup_go/controller"
	"os"
)

func main() {
	version := "0.9.3"
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
	c := controller.NewController("config.json", "dbs.txt")
	err = c.Controller()
	if err != nil {
		log.Panic(err)
	}
}
