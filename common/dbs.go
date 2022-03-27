package common

import (
	"bufio"
	"fmt"
	"os"
)

type ReadDbs interface {
	Read() ([]string, error)
}

type DBList struct {
	dblistFile string
	dbs        []string
}

func NewDBList(f string) *DBList {
	return &DBList{
		dblistFile: f,
	}
}

//读取数据库名
func (d *DBList) Read() (*[]string, error) {
	//打开文件
	dbsFile, err := os.OpenFile(d.dblistFile, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("创建数据库列表文件失败：%w", err)
	}
	defer dbsFile.Close()

	//按行读取文件并添加到切片
	bs := bufio.NewScanner(dbsFile)
	for {
		flag := bs.Scan()
		line := bs.Text()

		if flag == false {
			break
		}
		d.dbs = append(d.dbs, line)
	}
	return &d.dbs, nil
}
