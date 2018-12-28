package app

import (
	_ "github.com/Go-SQL-Driver/mysql"
	"github.com/go-ozzo/ozzo-dbx"
)

// LoadDB 创建DB
func LoadDB(dsn string) (*dbx.DB, error) {
	return dbx.MustOpen("mysql", dsn)
}
