package main

import (
	"log"
	"mysql_backup_go/controller"
	"os"
)

func main() {
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
