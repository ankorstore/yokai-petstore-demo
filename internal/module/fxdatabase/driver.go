package fxdatabase

import (
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type Driver struct {
	base  driver.Driver
	hooks []hook.Hook
}

func NewDriver(base driver.Driver, hooks []hook.Hook) *Driver {
	return &Driver{
		base:  base,
		hooks: hooks,
	}
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	connection, err := d.base.Open(name)
	if err != nil {
		return nil, err
	}

	return NewConnection(connection, d.hooks), nil
}

func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	connector, err := d.base.(driver.DriverContext).OpenConnector(name)
	if err != nil {
		return nil, err
	}

	return NewConnector(connector, d), nil
}

func (d *Driver) Hooks() []hook.Hook {
	return d.hooks
}
