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
	event := hook.NewHookEvent("Connection::Exec", query, args)

	ctx := c.applyBeforeHooks(context.Background(), event)

	engine, ok := c.base.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}

	res, err := engine.Exec(query, args)
	if err != nil {
		event.SetError(err)
	}

	if res != nil {
		lastInsertId, lastInsertIdErr := res.LastInsertId()
		if lastInsertIdErr == nil {
			event.SetLastInsertId(lastInsertId)
		}

		rowsAffected, rowsAffectedErr := res.RowsAffected()
		if rowsAffectedErr == nil {
			event.SetRowsAffected(rowsAffected)
		}
	}

	c.applyAfterHooks(ctx, event)

	return res, err
}

func (c *HookableConnection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	event := hook.NewHookEvent("Connection::ExecContext", query, args)

	ctx = c.applyBeforeHooks(ctx, event)

	engine, ok := c.base.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	res, err := engine.ExecContext(ctx, query, args)
	if err != nil {
		event.SetError(err)
	}

	if res != nil {
		lastInsertId, lastInsertIdErr := res.LastInsertId()
		if lastInsertIdErr == nil {
			event.SetLastInsertId(lastInsertId)
		}

		rowsAffected, rowsAffectedErr := res.RowsAffected()
		if rowsAffectedErr == nil {
			event.SetRowsAffected(rowsAffected)
		}
	}

	c.applyAfterHooks(ctx, event)

	return res, err
}

func (c *HookableConnection) Query(query string, args []driver.Value) (driver.Rows, error) {
	event := hook.NewHookEvent("Connection::Query", query, args)

	ctx := c.applyBeforeHooks(context.Background(), event)

	engine, ok := c.base.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}

	rows, err := engine.Query(query, args)
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return rows, err
}

func (c *HookableConnection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	event := hook.NewHookEvent("Connection::QueryContext", query, args)

	ctx = c.applyBeforeHooks(ctx, event)

	engine, ok := c.base.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	rows, err := engine.QueryContext(ctx, query, args)
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return rows, err
}

func (c *HookableConnection) Ping(ctx context.Context) error {
	event := hook.NewHookEvent("Connection::Ping", "ping", nil)

	ctx = c.applyBeforeHooks(ctx, event)

	engine, ok := c.base.(driver.Pinger)
	if !ok {
		return driver.ErrSkip
	}

	err := engine.Ping(ctx)
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return err
}

func (c *HookableConnection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	event := hook.NewHookEvent("Connection::PrepareContext", query, nil)

	ctx = c.applyBeforeHooks(ctx, event)

	if engine, ok := c.base.(driver.ConnPrepareContext); ok {
		stmt, err := engine.PrepareContext(ctx, query)
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewHookableStatement(stmt, ctx, query, c.hooks), err
	} else {
		stmt, err := c.base.Prepare(query)
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewHookableStatement(stmt, ctx, query, c.hooks), err
	}
}

func (c *HookableConnection) Prepare(query string) (driver.Stmt, error) {
	event := hook.NewHookEvent("Connection::Prepare", query, nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	stmt, err := c.base.Prepare(query)
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return NewHookableStatement(stmt, nil, query, c.hooks), err
}

func (c *HookableConnection) Begin() (driver.Tx, error) {
	event := hook.NewHookEvent("Connection::Begin", "", nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	tx, err := c.base.Begin()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return NewHookableTransaction(tx, ctx, c.hooks), err
}

func (c *HookableConnection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	event := hook.NewHookEvent("Connection::BeginTx", "", nil)

	ctx = c.applyBeforeHooks(ctx, event)

	if engine, ok := c.base.(driver.ConnBeginTx); ok {
		tx, err := engine.BeginTx(ctx, opts)
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewHookableTransaction(tx, ctx, c.hooks), err
	} else {
		tx, err := c.base.Begin()
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewHookableTransaction(tx, ctx, c.hooks), err
	}
}

func (c *HookableConnection) Close() error {
	event := hook.NewHookEvent("Connection::Close", "", nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	err := c.base.Close()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return err
}

func (c *HookableConnection) ResetSession(ctx context.Context) error {
	event := hook.NewHookEvent("Connection::ResetSession", "", nil)

	ctx = c.applyBeforeHooks(context.Background(), event)

	if engine, ok := c.base.(driver.SessionResetter); ok {
		err := engine.ResetSession(ctx)
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

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
