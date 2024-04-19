package fxdatabase

import (
	"github.com/ankorstore/yokai/log"
)

type MigrationLogger struct {
	logger *log.Logger
}

func NewMigrationLogger(logger *log.Logger) *MigrationLogger {
	return &MigrationLogger{
		logger: logger,
	}
}

func (l *MigrationLogger) Verbose() bool {
	return true
}

func (l *MigrationLogger) Printf(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}
