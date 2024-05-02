package fxdatabase

import (
	"context"
	"database/sql/driver"
)

type HookableConnector struct {
	dsn    string
	base   driver.Connector
	driver *HookableDriver
}

func NewHookableConnector(dsn string, base driver.Connector, driver *HookableDriver) *HookableConnector {
	return &HookableConnector{
		dsn:    dsn,
		base:   base,
		driver: driver,
	}
}

func (c *HookableConnector) Connect(ctx context.Context) (driver.Conn, error) {
	if c.base == nil {
		return c.driver.Open(c.dsn)
	}

	conn, err := c.base.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return NewHookableConnection(conn, c.driver.Configuration()), nil
}

func (c *HookableConnector) Driver() driver.Driver {
	return c.driver
}
