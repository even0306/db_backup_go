package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ReadConfig interface {
	Read() (*ConfigFile, error)
}

//基础配置
type ConfigFile struct {
	configFile string

	FILTER_METHOD    bool   `json:"FILTER_METHOD"`    //正向匹配 true，反向匹配 false
	REMOTE_BACKUP    bool   `json:"REMOTE_BACKUP"`    //开启向异机备份 true，关闭向异机备份 false
	SAVE_DAY         int    `json:"SAVE_DAY"`         //保存备份的天数
	MYSQL_EXEC_PATH  string `json:"MYSQL_EXEC_PATH"`  //mysql执行文件所在目录
	BACKUP_SAVE_PATH string `json:"BACKUP_SAVE_PATH"` //备份在本地保存的路径

	DB_HOST     string `json:"DB_HOST"`
	DB_PORT     int    `json:"DB_PORT"`
	DB_USER     string `json:"DB_USER"`
	DB_PASSWORD string `json:"DB_PASSWORD"`
	DB_LABEL    string `json:"DB_LABEL"` //标签，用于标记该备份来自哪个数据库

	REMOTE_HOST     string `json:"REMOTE_HOST"`
	REMOTE_PORT     int    `json:"REMOTE_PORT"`
	REMOTE_USER     string `json:"REMOTE_USER"`
	REMOTE_PASSWORD string `json:"REMOTE_PASSWORD"`
	REMOTE_PATH     string `json:"REMOTE_PATH"` //备份在异机保存的路径
}

func NewConfig(f string) *ConfigFile {
	return &ConfigFile{
		configFile:       f,
		FILTER_METHOD:    true,
		REMOTE_BACKUP:    false,
		SAVE_DAY:         7,
		MYSQL_EXEC_PATH:  "/usr/bin",
		BACKUP_SAVE_PATH: "",
		DB_HOST:          "127.0.0.1",
		DB_PORT:          3306,
		DB_USER:          "root",
		DB_PASSWORD:      "",
		DB_LABEL:         "",
		REMOTE_HOST:      "",
		REMOTE_PORT:      22,
		REMOTE_USER:      "",
		REMOTE_PASSWORD:  "",
		REMOTE_PATH:      "",
	}
}

//读取配置文件
func (c *ConfigFile) Read() (*ConfigFile, error) {
	//打开文件
	jsonFile, err := os.OpenFile(c.configFile, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("创建日志文件失败：%w", err)
	}
	defer jsonFile.Close()

	//读取配置文件
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%w", err)
	}

	//将配置文件内容传入结构体
	err = json.Unmarshal([]byte(byteValue), c)
	if err != nil {
		return nil, fmt.Errorf("配置文件内容转进程序失败：%w", err)
	}
	return c, nil
}
