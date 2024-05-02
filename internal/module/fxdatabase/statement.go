package fxdatabase

import (
	"context"
	"database/sql/driver"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
)

type HookableStatement struct {
	base          driver.Stmt
	context       context.Context
	query         string
	configuration *Configuration
}

func NewHookableStatement(base driver.Stmt, ctx context.Context, query string, configuration *Configuration) *HookableStatement {
	if ctx == nil {
		ctx = context.Background()
	}

	return &HookableStatement{
		base:          base,
		context:       ctx,
		query:         query,
		configuration: configuration,
	}
}

func (s *HookableStatement) Close() error {
	return s.base.Close()
}

func (s *HookableStatement) NumInput() int {
	return s.base.NumInput()
}

func (s *HookableStatement) Exec(args []driver.Value) (driver.Result, error) {
	event := hook.NewHookEvent(s.configuration.Driver(), "Statement::Exec", s.query, args)

	s.applyBeforeHooks(event)

	event.Start()
	res, err := s.base.Exec(args)
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

	s.applyAfterHooks(event)

	return res, err
}

func (s *HookableStatement) Query(args []driver.Value) (driver.Rows, error) {
	event := hook.NewHookEvent(s.configuration.Driver(), "Statement::Query", s.query, args)

	s.applyBeforeHooks(event)

	event.Start()
	rows, err := s.base.Query(args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	s.applyAfterHooks(event)

	return rows, err
}

func (s *HookableStatement) applyBeforeHooks(event *hook.HookEvent) {
	for _, h := range s.configuration.Hooks() {
		s.context = h.Before(s.context, event)
	}
}

func (s *HookableStatement) applyAfterHooks(event *hook.HookEvent) {
	for _, h := range s.configuration.Hooks() {
		h.After(s.context, event)
	}
}
