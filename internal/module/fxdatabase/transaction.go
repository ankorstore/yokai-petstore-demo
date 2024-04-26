package fxdatabase

import (
	"context"
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type HookableTransaction struct {
	base    driver.Tx
	context context.Context
	hooks   []hook.Hook
}

func NewHookableTransaction(base driver.Tx, ctx context.Context, hooks []hook.Hook) *HookableTransaction {
	if ctx == nil {
		ctx = context.Background()
	}

	return &HookableTransaction{
		base:    base,
		context: ctx,
		hooks:   hooks,
	}
}

func (t *HookableTransaction) Commit() error {
	event := hook.NewHookEvent("Transaction::Commit", "", nil)

	t.applyBeforeHooks(event)

	err := t.base.Commit()

	t.applyAfterHooks(event.SetError(err))

	return err
}

func (t *HookableTransaction) Rollback() error {
	event := hook.NewHookEvent("Transaction::Rollback", "", nil)

	t.applyBeforeHooks(event)

	err := t.base.Rollback()

	t.applyAfterHooks(event.SetError(err))

	return err
}

func (t *HookableTransaction) applyBeforeHooks(evt *hook.HookEvent) {
	for _, h := range t.hooks {
		if !t.checkHookExcluded(h, evt) {
			t.context = h.Before(t.context, evt)
		}
	}
}

func (t *HookableTransaction) applyAfterHooks(evt *hook.HookEvent) {
	for _, h := range t.hooks {
		if !t.checkHookExcluded(h, evt) {
			h.After(t.context, evt)
		}
	}
}

func (t *HookableTransaction) checkHookExcluded(h hook.Hook, evt *hook.HookEvent) bool {
	for _, exclusion := range h.Exclusions() {
		if evt.Name() == exclusion {
			return true
		}
	}

	return false
}
