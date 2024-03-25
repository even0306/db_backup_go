package shell

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetMysqlDBList(t *testing.T) {
	var info struct {
		DBUser     string
		DBHost     string
		DBPassword string
		DBPort     int
	}
	info.DBUser = "root"
	info.DBHost = "127.0.0.1"
	info.DBPassword = "123456"
	info.DBPort = 3306

	db, err := sql.Open("mysql", info.DBUser+":"+info.DBPassword+"@tcp("+info.DBHost+":"+fmt.Sprint(info.DBPort)+")/information_schema?charset=utf8")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	rows, err := db.Query("select schema_name from schemata")
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()
	if err != nil {
		t.Error(err)
	}
	var list []string
	for rows.Next() {
		var col string
		rows.Scan(&col)
		list = append(list, col)
	}
	for _, v := range list {
		fmt.Println(v)
	}
}
