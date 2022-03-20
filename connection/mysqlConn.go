package connection

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type connSocket interface {
	ConnSocket(db string, dbinfo string) (*sql.DB, error)
}

type Mysql struct {
	sql.DB
	db     string
	dbinfo string
}

func (ms *Mysql) ConnSocket(db string, dbinfo string) (*sql.DB, error) {
	ms.db = db
	ms.dbinfo = dbinfo
	conn, err := sql.Open(ms.db, ms.dbinfo)
	if err != nil {
		return nil, err
	}
	return conn, err
}
