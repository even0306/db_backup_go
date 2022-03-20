package main

import (
	"log"
	"mysql_backup_go/config"
	"mysql_backup_go/core"
)

func main() {
	logfile, logs := config.Init()
	log.SetOutput(logfile)
	err := core.Backup()
	if err != nil {
		logs.ErrorLogger.Panicln(err)
	}
}
