package fxdatabase

import "strings"

// Driver is an enum for the supported database drivers.
type Driver int

const (
	UnknownDriver Driver = iota
	Sqlite
	Mysql
	Postgres
)

// String returns a string representation of the [DriverName].
func (d Driver) String() string {
	switch d {
	case Sqlite:
		return "sqlite"
	case Mysql:
		return "mysql"
	case Postgres:
		return "postgres"
	default:
		return "unknown"
	}
}

// FetchDriver returns a [DriverName] for a given value.
func FetchDriver(driver string) Driver {
	switch strings.ToLower(driver) {
	case "sqlite":
		return Sqlite
	case "mysql":
		return Mysql
	case "postgres":
		return Postgres
	default:
		return UnknownDriver
	}
}

type MigrationDirection int

const (
	Unknown MigrationDirection = iota
	Up
	Down
)

func (d MigrationDirection) String() string {
	switch d {
	case Up:
		return "up"
	case Down:
		return "down"
	default:
		return "unknown"
	}
}

func FetchMigrationDirection(direction string) MigrationDirection {
	switch strings.ToLower(direction) {
	case "up":
		return Up
	case "down":
		return Down
	default:
		return Unknown
	}
}
