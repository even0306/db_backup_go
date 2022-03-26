package config

import (
	"bufio"
	"fmt"
	"os"
)

type readDbs interface {
	Read() ([]string, error)
}

type DBList struct {
	dblistFile string
}

func NewDBList(f string) *DBList {
	return &DBList{
		dblistFile: f,
	}
}

//读取数据库名
func (d *DBList) Read() ([]string, error) {
	//打开文件
	dbsFile, err := os.OpenFile(d.dblistFile, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("创建数据库列表文件失败：%v", err)
	}
	defer dbsFile.Close()

	//按行读取文件并添加到切片
	var dbs = []string{}
	line := bufio.NewReader(dbsFile)
	for {
		content, _, err := line.ReadLine()

		if err != nil && content != nil {
			return nil, fmt.Errorf("读取数据库列表行失败：%v", err)
		}
		if content == nil {
			return dbs, nil
		}
		dbs = append(dbs, string(content))
	}
}
