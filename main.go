package main

import (
	"db_backup_go/common"
	"db_backup_go/controller"
	"db_backup_go/logging"
	"flag"
	"os"
	"path/filepath"
)

//程序主入口
func main() {
	version := "1.1.9"
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
	if err == nil {
		_, err = os.Stat(exPath + "/dbs.txt")
		if err == nil {
			c := controller.NewController(exPath+"/config.json", exPath+"/dbs.txt")
			err = c.Controller()
			if err != nil {
				logging.Logger.Panic(err)
			}
		}
		if os.IsNotExist(err) {
			err := common.CreateDBS()
			if err != nil {
				logging.Logger.Panic(err)
			}
			main()
		}
	}
	if os.IsNotExist(err) {
		err := common.CreateConfig()
		if err != nil {
			logging.Logger.Panic(err)
		}
		main()
	}
}
