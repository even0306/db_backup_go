package main

import (
	"fmt"
	"mysql_backup_go/config"
)

func main() {
	a, b, c := config.ReadConfig()
	fmt.Println(a.BACKUP_SAVE_PATH)
	fmt.Println(b.DB_HOST)
	fmt.Println(c.REMOTE_PATH)
}
