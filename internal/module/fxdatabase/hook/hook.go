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
	name      string
	query     string
	arguments interface{}
	results   interface{}
	err       error
	timestamp time.Time
}

func NewHookEvent(name string, query string, arguments interface{}, results interface{}) *HookEvent {
	return &HookEvent{
		name:      name,
		query:     query,
		arguments: arguments,
		results:   results,
		timestamp: time.Now(),
	}
}

func (e *HookEvent) Name() string {
	return e.name
}

func (e *HookEvent) Query() string {
	return e.query
}

func (e *HookEvent) Arguments() interface{} {
	return e.arguments
}

func (e *HookEvent) SetResults(results interface{}) *HookEvent {
	e.results = results

	return e
}

func (e *HookEvent) Results() interface{} {
	return e.results
}

func (e *HookEvent) SetError(err error) *HookEvent {
	e.err = err

	return e
}

func (e *HookEvent) Error() error {
	return e.err
}

func (e *HookEvent) Timestamp() time.Time {
	return e.timestamp
}
