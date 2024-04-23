package fxsqlc

import (
	"database/sql"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"

	"go.uber.org/fx"
)

const ModuleName = "sqlc"

var FxSqlcModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxSqlcQueries,
	),
)

type FxSqlcQueriesParam struct {
	fx.In
	Db *sql.DB
}

func NewFxSqlcQueries(p FxSqlcQueriesParam) sqlc.Querier {
	return sqlc.New(p.Db)
}
