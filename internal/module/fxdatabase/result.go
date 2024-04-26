package fxdatabase

import (
	"context"
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type HookableResult struct {
	base    driver.Result
	context context.Context
	query   string
	hooks   []hook.Hook
}

func NewHookableResult(base driver.Result, ctx context.Context, query string, hooks []hook.Hook) *HookableResult {
	if ctx == nil {
		ctx = context.Background()
	}

	return &HookableResult{
		base:    base,
		context: ctx,
		query:   query,
		hooks:   hooks,
	}
}

func (r *HookableResult) LastInsertId() (int64, error) {
	event := hook.NewHookEvent("Result::LastInsertId", r.query, nil, nil)

	r.applyBeforeHooks(event)

	id, err := r.base.LastInsertId()

	r.applyAfterHooks(event.SetResults(map[string]interface{}{"lastInsertId": id}).SetError(err))

	return id, err
}

func (r *HookableResult) RowsAffected() (int64, error) {
	event := hook.NewHookEvent("Result::RowsAffected", r.query, nil, nil)

	r.applyBeforeHooks(event)

	rows, err := r.base.RowsAffected()

	r.applyAfterHooks(event.SetResults(map[string]interface{}{"rowsAffected": rows}).SetError(err))

	return rows, err
}

func (r *HookableResult) applyBeforeHooks(evt *hook.HookEvent) {
	for _, h := range r.hooks {
		if !r.checkHookExcluded(h, evt) {
			r.context = h.Before(r.context, evt)
		}
	}
}

func (r *HookableResult) applyAfterHooks(evt *hook.HookEvent) {
	for _, h := range r.hooks {
		if !r.checkHookExcluded(h, evt) {
			h.After(r.context, evt)
		}
	}
}

func (r *HookableResult) checkHookExcluded(h hook.Hook, evt *hook.HookEvent) bool {
	for _, exclusion := range h.Exclusions() {
		if evt.Name() == exclusion {
			return true
		}
	}

	return false
}
