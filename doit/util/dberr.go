package util

import (
	"database/sql"
	"github.com/Go-SQL-Driver/mysql"
)

//已存在
func IsDBDuplicatedErr(err error) bool {
	if dbErr, ok := err.(*mysql.MySQLError); ok {
		if dbErr.Number == 1062 {
			return true
		}
	}
	return false
}

//返回操作错误
func IsDBNotFound(err error) bool {
	return err == sql.ErrNoRows
}
