package main

import "fmt"

var (
	x byte = 1
	y byte = 2
	z byte = 2
)

func GetVersion() {
	fmt.Printf("db_backup_go\nversion: %v.%v.%v", x, y, z)
}
