package fxdatabase

import (
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type HookableDriver struct {
	name  string
	base  driver.Driver
	hooks []hook.Hook
}

func NewHookableDriver(name string, base driver.Driver, hooks []hook.Hook) *HookableDriver {
	return &HookableDriver{
		name:  name,
		base:  base,
		hooks: hooks,
	}
}

func (d *HookableDriver) Name() string {
	return d.name
}

func (d *HookableDriver) Open(dsn string) (driver.Conn, error) {
	connection, err := d.base.Open(dsn)
	if err != nil {
		return nil, err
	}

	return NewHookableConnection(connection, d.hooks), nil
}

func (d *HookableDriver) OpenConnector(dsn string) (driver.Connector, error) {
	if driverContext, ok := d.base.(driver.DriverContext); ok {
		connector, err := driverContext.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}

		return NewHookableConnector(dsn, connector, d), nil
	}

	return NewHookableConnector(dsn, nil, d), nil
}
