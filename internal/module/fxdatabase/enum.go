package fxdatabase

import "strings"

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
