package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

type ReadConfig interface {
	Read() error
}

//基础配置
type ConfigFile struct {
	configFile string

	DATABASETYPE *string `json:"database"`

	FILTER_METHOD    *bool   `json:"FILTER_METHOD"`    //正向匹配 true，反向匹配 false
	REMOTE_BACKUP    *bool   `json:"REMOTE_BACKUP"`    //开启向异机备份 true，关闭向异机备份 false
	SAVE_DAY         *int    `json:"SAVE_DAY"`         //保存备份的天数
	MYSQL_EXEC_PATH  *string `json:"MYSQL_EXEC_PATH"`  //mysql执行文件所在目录
	BACKUP_SAVE_PATH *string `json:"BACKUP_SAVE_PATH"` //备份在本地保存的路径

	DB_Version  *string `json:"DB_VERSION"`
	DB_HOST     *string `json:"DB_HOST"`
	DB_PORT     *int    `json:"DB_PORT"`
	DB_USER     *string `json:"DB_USER"`
	DB_PASSWORD *string `json:"DB_PASSWORD"`
	DB_LABEL    *string `json:"DB_LABEL"` //标签，用于标记该备份来自哪个数据库

	REMOTE_HOST     *string `json:"REMOTE_HOST"`
	REMOTE_PORT     *int    `json:"REMOTE_PORT"`
	REMOTE_USER     *string `json:"REMOTE_USER"`
	REMOTE_PASSWORD *string `json:"REMOTE_PASSWORD"`
	REMOTE_PATH     *string `json:"REMOTE_PATH"` //备份在异机保存的路径
}

// 初始化读取配置功能，传入string类型文件，返回*ConfigFile类型的结构体指针
func NewConfig(f string) *ConfigFile {
	return &ConfigFile{
		configFile:       f,
		DATABASETYPE:     new(string),
		FILTER_METHOD:    new(bool),
		REMOTE_BACKUP:    new(bool),
		SAVE_DAY:         new(int),
		MYSQL_EXEC_PATH:  new(string),
		BACKUP_SAVE_PATH: new(string),
		DB_Version:       new(string),
		DB_HOST:          new(string),
		DB_PORT:          new(int),
		DB_USER:          new(string),
		DB_PASSWORD:      new(string),
		DB_LABEL:         new(string),
		REMOTE_HOST:      new(string),
		REMOTE_PORT:      new(int),
		REMOTE_USER:      new(string),
		REMOTE_PASSWORD:  new(string),
		REMOTE_PATH:      new(string),
	}
}

//读取配置文件，返回error错误信息
func (c *ConfigFile) Read() error {
	//读取配置文件
	bf, err := ioutil.ReadFile(c.configFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败：%w", err)
	}

	//根据换行符转成string类型切片
	lines := strings.Split(string(bf), "\n")

	//遍历每一行，忽略行中带 # ; // 符号的数据，并去除行前后的无关空格，将其余数据存储起来
	var jsonValue []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "//") {
			continue
		}
		jsonValue = append(jsonValue, line)
	}

	//将新保存的数据的每一行末尾加上换行符，并追加到一起，成为独立字符串，方便后面转回字节流
	var build strings.Builder
	for _, v := range jsonValue {
		build.WriteString(v + "\n")
	}
	jsonByte := build.String()

	//将配置文件字节流转换为结构体
	err = json.Unmarshal([]byte(jsonByte), c)
	if err != nil {
		return fmt.Errorf("配置文件内容转进程序失败：%w", err)
	}

	k := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)
	for i := 1; i < v.NumField(); i++ {
		if k.Field(i).Tag.Get("json") == "" {
			log.Panicf("配置文件中缺少字段：%v", c.DATABASETYPE)
		}

		TypeAssertion(v)

	}
	return nil
}
