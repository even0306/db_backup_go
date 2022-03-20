package config

import (
	"bufio"
	"os"
)

type readDbs interface {
	Read(f string) ([]string, error)
}

type Databases struct {
}

func (d *Databases) Read(f string) ([]string, error) {
	dbsFile, err := os.OpenFile(f, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer dbsFile.Close()

	var dbs = []string{}
	line := bufio.NewReader(dbsFile)
	for {
		content, _, err := line.ReadLine()

		if err != nil && content != nil {
			return nil, err
		}
		if content == nil {
			return dbs, err
		}
		dbs = append(dbs, string(content))
	}
}
