package interfaces

import (
	"database/sql"
)

type SQLDriver interface {
	init()
	Close() error
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	Begin() (*sql.Tx, error)
	Commit() error
	Rollback() error
}
