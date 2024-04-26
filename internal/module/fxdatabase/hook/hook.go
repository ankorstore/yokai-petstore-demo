package hook

import (
	"context"
	"time"
)

type Hook interface {
	Exclusions() []string
	Before(context.Context, *HookEvent) context.Context
	After(context.Context, *HookEvent)
}

type HookEvent struct {
	name         string
	query        string
	arguments    any
	lastInsertId int64
	rowsAffected int64
	err          error
	timestamp    time.Time
}

func NewHookEvent(name string, query string, arguments interface{}) *HookEvent {
	return &HookEvent{
		name:      name,
		query:     query,
		arguments: arguments,
		timestamp: time.Now(),
	}
}

func (e *HookEvent) Name() string {
	return e.name
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

func (e *HookEvent) Timestamp() time.Time {
	return e.timestamp
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
