package hook

import (
	"context"
	"time"
)

type Hook interface {
	Before(context.Context, *HookEvent) context.Context
	After(context.Context, *HookEvent)
}

type HookEvent struct {
	driver       string
	operation    string
	query        string
	arguments    any
	lastInsertId int64
	rowsAffected int64
	err          error
	startedAt    time.Time
	stoppedAt    time.Time
}

func NewHookEvent(driver string, operation string, query string, arguments interface{}) *HookEvent {
	return &HookEvent{
		driver:    driver,
		operation: operation,
		query:     query,
		arguments: arguments,
	}
}

func (e *HookEvent) Start() *HookEvent {
	e.startedAt = time.Now()

	return e
}

func (e *HookEvent) Stop() *HookEvent {
	e.stoppedAt = time.Now()

	return e
}

func (e *HookEvent) Latency() time.Duration {
	return e.stoppedAt.Sub(e.startedAt)
}

func (e *HookEvent) Driver() string {
	return e.driver
}

func (e *HookEvent) Operation() string {
	return e.operation
}

func (e *HookEvent) Query() string {
	return e.query
}

func (e *HookEvent) Arguments() any {
	return e.arguments
}

func (e *HookEvent) LastInsertId() int64 {
	return e.lastInsertId
}
func (e *HookEvent) RowsAffected() int64 {
	return e.rowsAffected
}

func (e *HookEvent) Error() error {
	return e.err
}

func (e *HookEvent) SetLastInsertId(lastInsertId int64) *HookEvent {
	e.lastInsertId = lastInsertId

	return e
}

func (e *HookEvent) SetRowsAffected(rowsAffected int64) *HookEvent {
	e.rowsAffected = rowsAffected

	return e
}

func (e *HookEvent) SetError(err error) *HookEvent {
	e.err = err

	return e
}
