package main

import (
	"fmt"
	"log"
	"mysql_backup_go/controller"
	"os"
)

func main() {
	version := "0.9.2"
	for _, args := range os.Args {
		if len(os.Args) < 2 {
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
		} else {
			if args == "-v" || args == "--version" {
				fmt.Printf("mysql_backup\n" + version + "\n")
			} else {
				fmt.Printf("参数错误，可以使用以下参数：\n[-v|--version]\n")
			}
		}
	}
}
