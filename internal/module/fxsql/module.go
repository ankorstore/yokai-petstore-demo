package fxsql

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	yokaisql "github.com/ankorstore/yokai/sql"
	yokaisqllog "github.com/ankorstore/yokai/sql/hook/log"
	yokaisqltrace "github.com/ankorstore/yokai/sql/hook/trace"
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

const ModuleName = "sql"

var FxSQLModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxSQLDatabase,
	),
)

type FxSQLDatabaseParam struct {
	fx.In
	LifeCycle       fx.Lifecycle
	Config          *config.Config
	Logger          *log.Logger
	MetricsRegistry *prometheus.Registry
}

func NewFxSQLDatabase(p FxSQLDatabaseParam) (*sql.DB, error) {
	// hooks
	var driverHooks []yokaisql.Hook

	// trace hook
	if p.Config.GetBool("modules.database.trace.enabled") {
		driverHooks = append(
			driverHooks,
			yokaisqltrace.NewTraceHook(
				yokaisqltrace.WithArguments(p.Config.GetBool("modules.database.trace.arguments")),
				yokaisqltrace.WithExcludedOperations(yokaisql.FetchOperations(p.Config.GetStringSlice("modules.database.trace.exclude"))...),
			),
		)
	}

	// log hook
	if p.Config.GetBool("modules.database.log.enabled") {
		driverHooks = append(
			driverHooks,
			yokaisqllog.NewLogHook(
				yokaisqllog.WithLevel(log.FetchLogLevel(p.Config.GetString("modules.database.log.level"))),
				yokaisqllog.WithArguments(p.Config.GetBool("modules.database.log.arguments")),
				yokaisqllog.WithExcludedOperations(yokaisql.FetchOperations(p.Config.GetStringSlice("modules.database.log.exclude"))...),
			),
		)
	}

	// driver registration
	driverName, err := yokaisql.Register(p.Config.GetString("modules.database.driver"), driverHooks...)
	if err != nil {
		return nil, err
	}

	// database
	db, err := sql.Open(driverName, p.Config.GetString("modules.database.dsn"))
	if err != nil {
		return nil, err
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if p.Config.GetString("modules.database.driver") != "sqlite" {
				return db.Close()
			}

			return nil
		},
	})

	return db, nil
}

func RunFxSQLDatabaseMigration(command string, shutdown bool) fx.Option {
	return fx.Invoke(func(ctx context.Context, db *sql.DB, config *config.Config, logger *log.Logger, sd fx.Shutdowner) error {
		logger.Info().Str("command", command).Msg("starting database migration")

		err := goose.SetDialect(config.GetString("modules.database.driver"))
		if err != nil {
			logger.Error().Err(err).Str("command", command).Msg("database migration dialect error")

			return err
		}

		err = goose.RunContext(ctx, command, db, config.GetString("modules.database.migrations"))
		if err != nil {
			logger.Error().Err(err).Str("command", command).Msg("database migration error")

			return err
		}

		logger.Info().Str("command", command).Msg("database migration success")

		if shutdown {
			return sd.Shutdown()
		}

		return nil
	})
}
