package fxdatabase

import (
	"context"
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type HookableTransaction struct {
	base          driver.Tx
	context       context.Context
	configuration *Configuration
}

func NewHookableTransaction(base driver.Tx, ctx context.Context, configuration *Configuration) *HookableTransaction {
	if ctx == nil {
		ctx = context.Background()
	}

	return &HookableTransaction{
		base:          base,
		context:       ctx,
		configuration: configuration,
	}
}

func (t *HookableTransaction) Commit() error {
	event := hook.NewHookEvent(t.configuration.Driver(), "Transaction::Commit", "", nil)

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
	event := hook.NewHookEvent(t.configuration.Driver(), "Transaction::Rollback", "", nil)

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
	for _, h := range t.configuration.Hooks() {
		t.context = h.Before(t.context, event)
	}
}

func (t *HookableTransaction) applyAfterHooks(event *hook.HookEvent) {
	for _, h := range t.configuration.Hooks() {
		h.After(t.context, event)
	}
}
