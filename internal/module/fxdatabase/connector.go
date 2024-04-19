package fxdatabase

import (
	"context"
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type Connector struct {
	base   driver.Connector
	driver *Driver
	hooks  []hook.Hook
}

func NewConnector(base driver.Connector, driver *Driver) *Connector {
	return &Connector{
		base:   base,
		driver: driver,
	}
}

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.base.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return NewConnection(conn, c.driver.Hooks()), nil
}

func (c *Connector) Driver() driver.Driver {
	return c.driver
}
