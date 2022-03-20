package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Readconfig interface {
	Read(f string) (*Config, error)
}

//基础配置
type Config struct {
	FILTER_METHOD    bool   `json:"FILTER_METHOD"`    //正向匹配 true，反向匹配 false
	REMOTE_BACKUP    bool   `json:"REMOTE_BACKUP"`    //开启向异机备份 true，关闭向异机备份 false
	USE_KEY          bool   `json:"USE_KEY"`          //使用证书连接异机 true，使用密码连接异机 false
	SAVE_DAY         int    `json:"SAVE_DAY"`         //保存备份的天数
	MYSQL_EXEC_PATH  string `json:"MYSQL_EXEC_PATH"`  //mysql执行文件所在目录
	BACKUP_SAVE_PATH string `json:"BACKUP_SAVE_PATH"` //备份在本地保存的路径
	BACKUP_LOGS      string `json:"BACKUP_LOGS"`      //日志在本地的路径

	DB_HOST     string `json:"DB_HOST"`
	DB_PORT     string `json:"DB_PORT"`
	DB_USER     string `json:"DB_USER"`
	DB_PASSWORD string `json:"DB_PASSWORD"`
	DB_LABEL    string `json:"DB_LABEL"` //标签，用于标记该备份来自哪个数据库

	REMOTE_HOST string `json:"REMOTE_HOST"`
	REMOTE_PORT string `json:"REMOTE_PORT"`
	REMOTE_USER string `json:"REMOTE_USER"`
	REMOTE_KEY  string `json:"REMOTE_KEY"`
	REMOTE_PATH string `json:"REMOTE_PATH"` //备份在异机保存的路径
}

func (c *Config) Read(f string) (*Config, error) {
	jsonFile, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var config Config
	json.Unmarshal([]byte(byteValue), &config)
	return &config, err
}
