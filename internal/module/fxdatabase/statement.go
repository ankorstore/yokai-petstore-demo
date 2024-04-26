package fxdatabase

import (
	"context"
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type HookableStatement struct {
	base    driver.Stmt
	context context.Context
	query   string
	hooks   []hook.Hook
}

func NewHookableStatement(base driver.Stmt, ctx context.Context, query string, hooks []hook.Hook) *HookableStatement {
	if ctx == nil {
		ctx = context.Background()
	}

	return &HookableStatement{
		base:    base,
		context: ctx,
		query:   query,
		hooks:   hooks,
	}
}

func (s *HookableStatement) Close() error {
	return s.base.Close()
}

func (s *HookableStatement) NumInput() int {
	return s.base.NumInput()
}

func (s *HookableStatement) Exec(args []driver.Value) (driver.Result, error) {
	event := hook.NewHookEvent("Statement::Exec", s.query, args, nil)

	s.applyBeforeHooks(event)

	res, err := s.base.Exec(args)

	s.applyAfterHooks(event.SetResults(res).SetError(err))

	return NewHookableResult(res, s.context, s.query, s.hooks), err
}

func (s *HookableStatement) Query(args []driver.Value) (driver.Rows, error) {
	event := hook.NewHookEvent("Statement::Query", s.query, args, nil)

	s.applyBeforeHooks(event)

	rows, err := s.base.Query(args)

	s.applyAfterHooks(event.SetResults(rows).SetError(err))

	return rows, err
}

func (s *HookableStatement) applyBeforeHooks(evt *hook.HookEvent) {
	for _, h := range s.hooks {
		if !s.checkHookExcluded(h, evt) {
			s.context = h.Before(s.context, evt)
		}
	}
}

func (s *HookableStatement) applyAfterHooks(evt *hook.HookEvent) {
	for _, h := range s.hooks {
		if !s.checkHookExcluded(h, evt) {
			h.After(s.context, evt)
		}
	}
}

func (s *HookableStatement) checkHookExcluded(h hook.Hook, evt *hook.HookEvent) bool {
	for _, exclusion := range h.Exclusions() {
		if evt.Name() == exclusion {
			return true
		}
	}

	return false
}
