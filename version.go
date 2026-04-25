package main

import "fmt"

var (
	x byte = 2
	y byte = 0
	z byte = 0
)

func GetVersion() {
	fmt.Printf("db_backup_go\nversion: %v.%v.%v\n", x, y, z)
}
