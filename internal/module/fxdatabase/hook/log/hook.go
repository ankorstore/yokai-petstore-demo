package log

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
	"github.com/ankorstore/yokai/log"
)

type LogHook struct{}

func NewLogHook() *LogHook {
	return &LogHook{}
}

func (h *LogHook) ExcludedOperations() []string {
	return []string{
		"Connection::Ping",
		"Connection::ResetSession",
	}
}

func (h *LogHook) Before(ctx context.Context, _ *hook.HookEvent) context.Context {
	return ctx
}

func (h *LogHook) After(ctx context.Context, event *hook.HookEvent) {
	logger := log.CtxLogger(ctx)

	loggerEvent := logger.Info()
	if event.Error() != nil {
		if !errors.Is(event.Error(), driver.ErrSkip) {
			loggerEvent = logger.Error().Err(event.Error())
		}
	}

	loggerEvent.
		Str("operation", event.Operation()).
		Str("latency", event.Latency().String())

	if event.Query() != "" {
		loggerEvent.Str("query", event.Query())
	}

	if event.Arguments() != nil {
		loggerEvent.Interface("arguments", event.Arguments())
	}

	loggerEvent.
		Int64("lastInsertId", event.LastInsertId()).
		Int64("rowsAffected", event.RowsAffected()).
		Msg("sql logger")
}
