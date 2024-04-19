package log

import (
	"context"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
	"github.com/ankorstore/yokai/log"
)

type LogHook struct{}

func NewLogHook() *LogHook {
	return &LogHook{}
}

func (h *LogHook) Exclusions() []string {
	return []string{
		"Ping",
		"ResetSession",
		"PrepareContext",
	}
}

func (h *LogHook) Before(ctx context.Context, event *hook.HookEvent) context.Context {

	return ctx
}

func (h *LogHook) After(ctx context.Context, event *hook.HookEvent) {
	if event.Name() != "Ping" && event.Name() != "ResetSession" {
		latency := time.Since(event.Timestamp())

		logger := log.CtxLogger(ctx)

		loggerEvent := logger.Info()
		if event.Error() != nil {
			if !errors.Is(event.Error(), driver.ErrSkip) {
				loggerEvent = logger.Error().Err(event.Error())
			}
		}

		loggerEvent.Str("operation", event.Name())
		loggerEvent.Str("query", event.Query())
		loggerEvent.Str("latency", latency.String())
		loggerEvent.Interface("args", event.Args())
		loggerEvent.Msg("sql logger")
	}
}
