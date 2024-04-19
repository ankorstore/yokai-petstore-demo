package trace

import (
	"context"
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

func (h *TraceHook) Exclusions() []string {
	return []string{
		"Ping",
		"ResetSession",
		"PrepareContext",
	}
}

func (h *TraceHook) Before(ctx context.Context, event *hook.HookEvent) context.Context {
	ctx, _ = trace.CtxTracerProvider(ctx).Tracer("yokai-sql").Start(
		ctx,
		event.Name(),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(h.eventAttributes(event)...),
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
		span.RecordError(event.Error())
		code = codes.Error
	}
	span.SetStatus(code, code.String())
	span.End()
}

func (h *TraceHook) eventAttributes(event *hook.HookEvent) []attribute.KeyValue {
	attributes := []attribute.KeyValue{
		semconv.DBStatementKey.String(event.Query()),
	}

	if event.Args() != nil {
		attributes = append(
			attributes,
			attribute.String("db.statement.args", fmt.Sprintf("%+v", event.Args())),
		)
	}

	return attributes
}
