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
	engine, ok := c.base.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}

	event := hook.NewHookEvent("Connection::Exec", query, args)

	ctx := c.applyBeforeHooks(context.Background(), event)

	event.Start()
	res, err := engine.Exec(query, args)
	event.Stop()
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
	engine, ok := c.base.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	event := hook.NewHookEvent("Connection::ExecContext", query, args)

	ctx = c.applyBeforeHooks(ctx, event)

	event.Start()
	res, err := engine.ExecContext(ctx, query, args)
	event.Stop()
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
	engine, ok := c.base.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}

	event := hook.NewHookEvent("Connection::Query", query, args)

	ctx := c.applyBeforeHooks(context.Background(), event)

	event.Start()
	rows, err := engine.Query(query, args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return rows, err
}

func (c *HookableConnection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	engine, ok := c.base.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	event := hook.NewHookEvent("Connection::QueryContext", query, args)

	ctx = c.applyBeforeHooks(ctx, event)

	event.Start()
	rows, err := engine.QueryContext(ctx, query, args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return rows, err
}

func (c *HookableConnection) Ping(ctx context.Context) error {
	engine, ok := c.base.(driver.Pinger)
	if !ok {
		return driver.ErrSkip
	}

	event := hook.NewHookEvent("Connection::Ping", "ping", nil)

	ctx = c.applyBeforeHooks(ctx, event)

	event.Start()
	err := engine.Ping(ctx)
	event.Stop()
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
		event.Start()
		stmt, err := engine.PrepareContext(ctx, query)
		event.Stop()
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewHookableStatement(stmt, ctx, query, c.hooks), err
	} else {
		event.Start()
		stmt, err := c.base.Prepare(query)
		event.Stop()
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

	event.Start()
	stmt, err := c.base.Prepare(query)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return NewHookableStatement(stmt, nil, query, c.hooks), err
}

func (c *HookableConnection) Begin() (driver.Tx, error) {
	event := hook.NewHookEvent("Connection::Begin", "", nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	event.Start()
	tx, err := c.base.Begin()
	event.Stop()
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
		event.Start()
		tx, err := engine.BeginTx(ctx, opts)
		event.Stop()
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewHookableTransaction(tx, ctx, c.hooks), err
	} else {
		event.Start()
		tx, err := c.base.Begin()
		event.Stop()
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

	event.Start()
	err := c.base.Close()
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return err
}

func (c *HookableConnection) ResetSession(ctx context.Context) error {
	engine, ok := c.base.(driver.SessionResetter)
	if !ok {
		return driver.ErrSkip
	}

	event := hook.NewHookEvent("Connection::ResetSession", "", nil)

	ctx = c.applyBeforeHooks(context.Background(), event)

	event.Start()
	err := engine.ResetSession(ctx)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return err

}

func (c *HookableConnection) applyBeforeHooks(ctx context.Context, event *hook.HookEvent) context.Context {
	for _, h := range c.hooks {
		ctx = h.Before(ctx, event)
	}

	return ctx
}

func (c *HookableConnection) applyAfterHooks(ctx context.Context, event *hook.HookEvent) {
	for _, h := range c.hooks {
		h.After(ctx, event)
	}
}
