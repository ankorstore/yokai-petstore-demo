package mysql

import "github.com/go-sql-driver/mysql"

func NewBaseDriver() *mysql.MySQLDriver {
	return &mysql.MySQLDriver{}
}
