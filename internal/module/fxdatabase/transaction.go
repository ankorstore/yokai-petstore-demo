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

	event.Start()
	err := t.base.Commit()
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	t.applyAfterHooks(event)

	return err
}

func (t *HookableTransaction) Rollback() error {
	event := hook.NewHookEvent("Transaction::Rollback", "", nil)

	t.applyBeforeHooks(event)

	event.Start()
	err := t.base.Rollback()
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	t.applyAfterHooks(event)

	return err
}

func (t *HookableTransaction) applyBeforeHooks(event *hook.HookEvent) {
	for _, h := range t.hooks {
		if !t.checkHookExcluded(h, event) {
			t.context = h.Before(t.context, event)
		}
	}
}

func (t *HookableTransaction) applyAfterHooks(event *hook.HookEvent) {
	for _, h := range t.hooks {
		if !t.checkHookExcluded(h, event) {
			h.After(t.context, event)
		}
	}
}

func (t *HookableTransaction) checkHookExcluded(h hook.Hook, event *hook.HookEvent) bool {
	for _, operation := range h.ExcludedOperations() {
		if event.Operation() == operation {
			return true
		}
	}

	return false
}
