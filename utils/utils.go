package utils

import (
	"fmt"
	"os/exec"
)

type check interface {
	CheckOS() (string, error)
}

type OSType struct {
	checkFunc map[string]string
}

func NewUtils() *OSType {
	return &OSType{
		checkFunc: map[string]string{},
	}
}

func (t *OSType) CheckOS() (string, error) {
	t.checkFunc = map[string]string{"linux": "uname", "windows": "systeminfo"}
	for k, v := range t.checkFunc {
		cmd := exec.Command(v)
		err := cmd.Start()
		if err == nil {
			return k, nil
		}
	}
	return "", fmt.Errorf("未知远端系统，发送失败")
}
