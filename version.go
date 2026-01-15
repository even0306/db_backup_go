package main

import "fmt"

var (
	x byte = 1
	y byte = 3
	z byte = 2
)

func GetVersion() {
	fmt.Printf("db_backup_go\nversion: %v.%v.%v\n", x, y, z)
}
