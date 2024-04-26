package fxdatabase

import (
	"context"
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type HookableConnection struct {
	base  driver.Conn
	hooks []hook.Hook
}

func NewHookableConnection(base driver.Conn, hooks []hook.Hook) *HookableConnection {
	return &HookableConnection{
		base:  base,
		hooks: hooks,
	}
}

func (c *HookableConnection) Exec(query string, args []driver.Value) (driver.Result, error) {
	event := hook.NewHookEvent("Connection::Exec", query, args, nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	engine, ok := c.base.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}

	res, err := engine.Exec(query, args)

	c.applyAfterHooks(ctx, event.SetResults(res).SetError(err))

	return NewHookableResult(res, ctx, query, c.hooks), err
}

func (c *HookableConnection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	event := hook.NewHookEvent("Connection::ExecContext", query, args, nil)

	ctx = c.applyBeforeHooks(ctx, event)

	engine, ok := c.base.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	res, err := engine.ExecContext(ctx, query, args)

	c.applyAfterHooks(ctx, event.SetResults(res).SetError(err))

	return NewHookableResult(res, ctx, query, c.hooks), err
}

func (c *HookableConnection) Query(query string, args []driver.Value) (driver.Rows, error) {
	event := hook.NewHookEvent("Connection::Query", query, args, nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	engine, ok := c.base.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}

	rows, err := engine.Query(query, args)

	c.applyAfterHooks(ctx, event.SetResults(rows).SetError(err))

	return rows, err
}

func (c *HookableConnection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	event := hook.NewHookEvent("Connection::QueryContext", query, args, nil)

	ctx = c.applyBeforeHooks(ctx, event)

	engine, ok := c.base.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	rows, err := engine.QueryContext(ctx, query, args)

	c.applyAfterHooks(ctx, event.SetResults(rows).SetError(err))

	return rows, err
}

func (c *HookableConnection) Ping(ctx context.Context) error {
	event := hook.NewHookEvent("Connection::Ping", "ping", nil, nil)

	ctx = c.applyBeforeHooks(ctx, event)

	engine, ok := c.base.(driver.Pinger)
	if !ok {
		return driver.ErrSkip
	}

	err := engine.Ping(ctx)

	c.applyAfterHooks(ctx, event.SetError(err))

	return err
}

func (c *HookableConnection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	event := hook.NewHookEvent("Connection::PrepareContext", query, nil, nil)

	ctx = c.applyBeforeHooks(ctx, event)

	if engine, ok := c.base.(driver.ConnPrepareContext); ok {
		stmt, err := engine.PrepareContext(ctx, query)

		c.applyAfterHooks(ctx, event.SetError(err))

		return NewHookableStatement(stmt, ctx, query, c.hooks), err
	} else {
		stmt, err := c.base.Prepare(query)

		c.applyAfterHooks(ctx, event.SetError(err))

		return NewHookableStatement(stmt, ctx, query, c.hooks), err
	}
}

func (c *HookableConnection) Prepare(query string) (driver.Stmt, error) {
	event := hook.NewHookEvent("Connection::Prepare", query, nil, nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	stmt, err := c.base.Prepare(query)

	c.applyAfterHooks(ctx, event.SetError(err))

	return NewHookableStatement(stmt, nil, query, c.hooks), err
}

func (c *HookableConnection) Begin() (driver.Tx, error) {
	event := hook.NewHookEvent("Connection::Begin", "begin", nil, nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	tx, err := c.base.Begin()

	c.applyAfterHooks(ctx, event.SetError(err))

	return tx, err
}

func (c *HookableConnection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	event := hook.NewHookEvent("Connection::BeginTx", "begin tx", nil, nil)

	ctx = c.applyBeforeHooks(ctx, event)

	if engine, ok := c.base.(driver.ConnBeginTx); ok {
		tx, err := engine.BeginTx(ctx, opts)

		c.applyAfterHooks(ctx, event.SetError(err))

		return tx, err
	} else {
		tx, err := c.base.Begin()

		c.applyAfterHooks(ctx, event.SetError(err))

		return tx, err
	}
}

func (c *HookableConnection) Close() error {
	event := hook.NewHookEvent("Connection::Close", "close", nil, nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	err := c.base.Close()

	c.applyAfterHooks(ctx, event.SetError(err))

	return err
}

func (c *HookableConnection) ResetSession(ctx context.Context) error {
	event := hook.NewHookEvent("Connection::ResetSession", "reset", nil, nil)

	ctx = c.applyBeforeHooks(context.Background(), event)

	if engine, ok := c.base.(driver.SessionResetter); ok {
		err := engine.ResetSession(ctx)

		c.applyAfterHooks(ctx, event.SetError(err))

		return err
	}

	c.applyAfterHooks(ctx, event)

	return nil
}

func (c *HookableConnection) applyBeforeHooks(ctx context.Context, evt *hook.HookEvent) context.Context {
	for _, h := range c.hooks {
		if !c.checkHookExcluded(h, evt) {
			ctx = h.Before(ctx, evt)
		}
	}

	return ctx
}

func (c *HookableConnection) applyAfterHooks(ctx context.Context, evt *hook.HookEvent) {
	for _, h := range c.hooks {
		if !c.checkHookExcluded(h, evt) {
			h.After(ctx, evt)
		}
	}
}

func (c *HookableConnection) checkHookExcluded(h hook.Hook, evt *hook.HookEvent) bool {
	for _, exclusion := range h.Exclusions() {
		if evt.Name() == exclusion {
			return true
		}
	}

	return false
}
