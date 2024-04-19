package fxdatabase

import (
	"context"
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type Connection struct {
	base  driver.Conn
	hooks []hook.Hook
}

func NewConnection(base driver.Conn, hooks []hook.Hook) *Connection {
	return &Connection{
		base:  base,
		hooks: hooks,
	}
}

func (c *Connection) Exec(query string, args []driver.Value) (res driver.Result, err error) {
	event := hook.NewHookEvent("Exec", query, args)

	ctx := c.applyBeforeHooks(context.Background(), event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	engine, ok := c.base.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}

	return engine.Exec(query, args)
}

func (c *Connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (res driver.Result, err error) {
	event := hook.NewHookEvent("ExecContext", query, args)

	ctx = c.applyBeforeHooks(ctx, event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	engine, ok := c.base.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	return engine.ExecContext(ctx, query, args)
}

func (c *Connection) Query(query string, args []driver.Value) (rows driver.Rows, err error) {
	event := hook.NewHookEvent("Query", query, args)

	ctx := c.applyBeforeHooks(context.Background(), event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	engine, ok := c.base.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}

	return engine.Query(query, args)
}

func (c *Connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	event := hook.NewHookEvent("QueryContext", query, args)

	ctx = c.applyBeforeHooks(ctx, event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	engine, ok := c.base.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	return engine.QueryContext(ctx, query, args)
}

func (c *Connection) Ping(ctx context.Context) (err error) {
	event := hook.NewHookEvent("Ping", "ping", nil)

	ctx = c.applyBeforeHooks(ctx, event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	engine, ok := c.base.(driver.Pinger)
	if !ok {
		return driver.ErrSkip
	}

	return engine.Ping(ctx)
}

func (c *Connection) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	event := hook.NewHookEvent("PrepareContext", query, nil)

	ctx = c.applyBeforeHooks(ctx, event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	if prepare, ok := c.base.(driver.ConnPrepareContext); ok {
		if stmt, err = prepare.PrepareContext(ctx, query); err != nil {
			return nil, err
		}
	} else {
		if stmt, err = c.base.Prepare(query); err != nil {
			return nil, err
		}
	}

	return stmt, err
}

func (c *Connection) Prepare(query string) (stmt driver.Stmt, err error) {
	event := hook.NewHookEvent("Prepare", query, nil)

	ctx := c.applyBeforeHooks(context.Background(), event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	return c.base.Prepare(query)
}

func (c *Connection) Begin() (tx driver.Tx, err error) {
	event := hook.NewHookEvent("Begin", "begin", nil)

	ctx := c.applyBeforeHooks(context.Background(), event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	return c.base.Begin()
}

func (c *Connection) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	event := hook.NewHookEvent("BeginTx", "begin tx", nil)

	ctx = c.applyBeforeHooks(ctx, event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	if beginTx, ok := c.base.(driver.ConnBeginTx); ok {
		if tx, err = beginTx.BeginTx(ctx, opts); err != nil {
			return nil, err
		}
	} else {
		if tx, err = c.base.Begin(); err != nil { // nolint
			return nil, err
		}
	}

	return tx, err
}

func (c *Connection) Close() (err error) {
	event := hook.NewHookEvent("Close", "close", nil)

	ctx := c.applyBeforeHooks(context.Background(), event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	return c.base.Close()
}

func (c *Connection) ResetSession(ctx context.Context) (err error) {
	event := hook.NewHookEvent("ResetSession", "reset", nil)

	ctx = c.applyBeforeHooks(context.Background(), event)
	defer func() {
		event.SetError(err)
		c.applyAfterHooks(ctx, event)
	}()

	if cr, ok := c.base.(driver.SessionResetter); ok {
		return cr.ResetSession(ctx)
	}

	return nil
}

func (c *Connection) applyBeforeHooks(ctx context.Context, evt *hook.HookEvent) context.Context {
	for _, h := range c.hooks {
		if !c.checkHookExcluded(h, evt) {
			ctx = h.Before(ctx, evt)
		}
	}

	return ctx
}

func (c *Connection) applyAfterHooks(ctx context.Context, evt *hook.HookEvent) {
	for _, h := range c.hooks {
		if !c.checkHookExcluded(h, evt) {
			h.After(ctx, evt)
		}
	}
}

func (c *Connection) checkHookExcluded(h hook.Hook, evt *hook.HookEvent) bool {
	for _, exclusion := range h.Exclusions() {
		if evt.Name() == exclusion {
			return true
		}
	}

	return false
}
