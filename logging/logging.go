package logging

import (
	"io"
	"log"
	"os"
)

var Logger *log.Logger

//初始化日志，创建并打开 server.log 文件，创建名为 Logger 的 *logging.Logger.Logger
func NewLogger(file string) {
	//初始化日志
	logfile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		Logger.Panicf("打开日志文件失败：%v", err)
	}
	stdout := os.Stdout
	Logger = log.New(io.MultiWriter(stdout, logfile), "", log.Lshortfile|log.LstdFlags)
}
