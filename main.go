package main

import (
	"log"
	"mysql_backup_go/controller"
)

func main() {
	c := controller.NewController()
	err := c.Controller()
	if err != nil {
		log.Panic(err)
	}
}
