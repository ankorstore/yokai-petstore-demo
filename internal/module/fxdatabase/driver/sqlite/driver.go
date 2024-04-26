package sqlite

import (
	"github.com/mattn/go-sqlite3"
)

func NewBaseDriver() *sqlite3.SQLiteDriver {
	return &sqlite3.SQLiteDriver{}
}
