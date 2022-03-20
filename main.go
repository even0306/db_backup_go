package main

import (
	"log"
	"mysql_backup_go/config"
	"mysql_backup_go/core"
)

func main() {
	logger := config.Logger{}
	logFile, logs := logger.SetLogConfig("server.log")
	log.SetOutput(logFile)
	err := core.Backup()
	if err != nil {
		logs.ErrorLogger.Panicln(err)
	}
}
