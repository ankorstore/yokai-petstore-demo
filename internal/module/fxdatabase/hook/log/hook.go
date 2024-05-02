package log

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
)

type LogHook struct {
	options Options
}

func NewLogHook(options ...LogHookOption) *LogHook {
	appliedOpts := DefaultLogHookOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	return &LogHook{
		options: appliedOpts,
	}
}

func (h *LogHook) ExcludedOperations() []string {
	return h.options.ExcludedOperations
}

func (h *LogHook) Before(ctx context.Context, _ *hook.HookEvent) context.Context {
	return ctx
}

func (h *LogHook) After(ctx context.Context, event *hook.HookEvent) {
	if hook.Contains(h.options.ExcludedOperations, event.Operation()) {
		return
	}

	logger := log.CtxLogger(ctx)

	var loggerEvent *zerolog.Event
	switch h.options.Level {
	case zerolog.DebugLevel:
		loggerEvent = logger.Debug()
	case zerolog.WarnLevel:
		loggerEvent = logger.Warn()
	default:
		loggerEvent = logger.Info()
	}

	if event.Error() != nil {
		if !errors.Is(event.Error(), driver.ErrSkip) {
			loggerEvent = logger.Error().Err(event.Error())
		}
	}

	loggerEvent.
		Str("driver", event.Driver()).
		Str("operation", event.Operation()).
		Str("latency", event.Latency().String())

	if event.Query() != "" {
		loggerEvent.Str("query", event.Query())
	}

	if h.options.Arguments && event.Arguments() != nil {
		loggerEvent.Interface("arguments", event.Arguments())
	}

	loggerEvent.
		Int64("lastInsertId", event.LastInsertId()).
		Int64("rowsAffected", event.RowsAffected()).
		Msg("sql logger")
}
