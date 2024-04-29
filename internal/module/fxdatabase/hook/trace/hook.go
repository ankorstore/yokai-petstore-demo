package trace

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
	"github.com/ankorstore/yokai/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type TraceHook struct{}

func NewTraceHook() *TraceHook {
	return &TraceHook{}
}

func (h *TraceHook) ExcludedOperations() []string {
	return []string{
		"Connection::Ping",
		"Connection::ResetSession",
	}
}

func (h *TraceHook) Before(ctx context.Context, event *hook.HookEvent) context.Context {
	var attributes []attribute.KeyValue

	if event.Query() != "" {
		attributes = append(
			attributes,
			semconv.DBStatementKey.String(event.Query()),
		)
	}

	if event.Arguments() != nil {
		attributes = append(
			attributes,
			attribute.String("db.statement.arguments", fmt.Sprintf("%+v", event.Arguments())),
		)
	}

	ctx, _ = trace.CtxTracerProvider(ctx).Tracer("yokai-sql").Start(
		ctx,
		event.Operation(),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(attributes...),
	)

	return ctx
}

func (h *TraceHook) After(ctx context.Context, event *hook.HookEvent) {
	span := oteltrace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	code := codes.Ok
	if event.Error() != nil {
		if !errors.Is(event.Error(), driver.ErrSkip) {
			code = codes.Error
			span.RecordError(event.Error())
		}
	}
	span.SetStatus(code, code.String())

	span.SetAttributes(
		attribute.String("db.latency", event.Latency().String()),
		attribute.Int64("db.lastInsertId", event.LastInsertId()),
		attribute.Int64("db.rowsAffected", event.RowsAffected()),
	)

	span.End()
}
