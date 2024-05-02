package fxdatabase

import (
	"database/sql/driver"
)

type HookableDriver struct {
	base          driver.Driver
	configuration *Configuration
}

func NewHookableDriver(base driver.Driver, configuration *Configuration) *HookableDriver {
	return &HookableDriver{
		base:          base,
		configuration: configuration,
	}
}

func (d *HookableDriver) Name() string {
	return d.configuration.Driver()
}

func (d *HookableDriver) Configuration() *Configuration {
	return d.configuration
}

func (d *HookableDriver) Open(dsn string) (driver.Conn, error) {
	connection, err := d.base.Open(dsn)
	if err != nil {
		return nil, err
	}

	return NewHookableConnection(connection, d.configuration), nil
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
