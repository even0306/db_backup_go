package utils

import (
	"log"
	"mysql_backup_go/config"
	"os/exec"
)

type check interface {
	CheckOS() string
}

type Type struct {
	osType map[string]string
}

func (t *Type) CheckOS() string {
	l := config.Logger{}
	logfile, logs := l.SetLogConfig("server.log")
	log.SetOutput(logfile)
	t.osType = map[string]string{"linux": "uname", "windows": "systeminfo"}
	for k, v := range t.osType {
		cmd := exec.Command(v)
		err := cmd.Start()
		if err == nil {
			return k
		}
	}
	logs.ErrorLogger.Println("未知远端系统，发送失败")
	return ""
}
