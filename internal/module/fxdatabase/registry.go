package fxdatabase

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/driver/mysql"
	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/driver/postgres"
	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/driver/sqlite"
	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

const DriverNamePrefix = "yokai"

var RegisteredDrivers *DriverRegistry

func init() {
	RegisteredDrivers = NewDriverRegistry()
}

type DriverRegistry struct {
	drivers map[string]driver.Driver
	mutex   sync.RWMutex
}

func NewDriverRegistry() *DriverRegistry {
	return &DriverRegistry{
		drivers: make(map[string]driver.Driver),
	}
}

func (r *DriverRegistry) Add(driverName string, driver driver.Driver) *DriverRegistry {
	if r.Has(driverName) {
		return r
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	r.drivers[driverName] = driver

	return r
}

func (r *DriverRegistry) Has(driverName string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, ok := r.drivers[driverName]

	return ok
}

func (r *DriverRegistry) Get(driverName string) (driver.Driver, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	registeredDriver, ok := r.drivers[driverName]
	if !ok {
		return nil, fmt.Errorf("cannot find driver %s in driver registry", driverName)
	}

	return registeredDriver, nil
}

func RegisterDriver(driverName string, hooks ...hook.Hook) (string, error) {
	driverRegistrationName := fmt.Sprintf("%s-%s", DriverNamePrefix, driverName)

	if RegisteredDrivers.Has(driverRegistrationName) {
		return driverRegistrationName, nil
	}

	var hookableDriver *HookableDriver
	switch FetchDriver(driverName) {
	case Mysql:
		hookableDriver = NewHookableDriver(mysql.NewBaseDriver(), NewConfiguration(driverName, hooks...))
	case Postgres:
		hookableDriver = NewHookableDriver(postgres.NewBaseDriver(), NewConfiguration(driverName, hooks...))
	case Sqlite:
		hookableDriver = NewHookableDriver(sqlite.NewBaseDriver(), NewConfiguration(driverName, hooks...))
	default:
		return "", fmt.Errorf("unsupported driver")
	}

	sql.Register(driverRegistrationName, hookableDriver)

	RegisteredDrivers.Add(driverRegistrationName, hookableDriver)

	return driverRegistrationName, nil
}
