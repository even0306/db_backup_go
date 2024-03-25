package shell

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

// 使用postgresql客户端查看postgresql数据库现有的库，返回*[]string的数据库列表切片指针
func TestGetPostgresqlDBList(t *testing.T) {
	var info struct {
		DBUser     string
		DBHost     string
		DBPassword string
		DBPort     int
	}
	info.DBUser = "postgres"
	info.DBHost = "127.0.0.1"
	info.DBPassword = "123456"
	info.DBPort = 5432

	db, err := sql.Open("postgres", "host="+info.DBHost+" port="+fmt.Sprint(info.DBPort)+" user="+info.DBUser+" password="+info.DBPassword+" dbname=postgres"+" sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	rows, err := db.Query("select datname from pg_catalog.pg_database")
	if err != nil {
		t.Error(err)
	}

	var list []string
	for rows.Next() {
		var col string
		err = rows.Scan(&col)
		if err != nil {
			t.Error(err)
		}
		list = append(list, col)
	}

	for _, v := range list {
		fmt.Println(v)
	}

	// cmd := exec.Command(info.ExecPath+"/psql", fmt.Sprintf("host=%s port=%v user=%s password=%s", info.DBHost, info.DBPort, info.DBUser, info.DBPassword), "-c", "SELECT datname FROM pg_database;")
	// var stdout bytes.Buffer
	// var stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	// err := cmd.Run()
	// if err != nil {
	// 	return nil, fmt.Errorf("数据库列表查询失败：%w:%v", err, stderr.String())
	// }
	// out := stdout.String()
	// list := strings.Split(string(out), "\n")
	// list = list[2 : len(list)-3]
	// return &list, nil
}
