package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type dbs interface{}

func ReadDbs() {
	file, logs := Init()
	log.SetOutput(file)
	dbsFile, err := os.Open("dbs.txt")
	if err != nil {
		logs.ErrorLogger.Panicln(err)
	}
	byteValue, err := ioutil.ReadAll(dbsFile)
	if err != nil {
		logs.ErrorLogger.Panicln(err)
	}
	fmt.Print(byteValue)
}
