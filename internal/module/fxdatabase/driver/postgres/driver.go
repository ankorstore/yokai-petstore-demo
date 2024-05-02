package postgres

import (
	"github.com/lib/pq"
)

func NewBaseDriver() *pq.Driver {
	return &pq.Driver{}
}
