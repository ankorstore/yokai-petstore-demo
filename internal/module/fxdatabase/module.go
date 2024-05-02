package fxdatabase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook"
	loghook "github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook/log"
	tracehook "github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase/hook/trace"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	postgresmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	sqlite3migrate "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

const ModuleName = "database"

var FxDatabaseModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxDatabase,
		NewFxDatabaseMigrator,
	),
)

type FxDatabaseParam struct {
	fx.In
	LifeCycle       fx.Lifecycle
	Config          *config.Config
	Logger          *log.Logger
	MetricsRegistry *prometheus.Registry
}

func NewFxDatabase(p FxDatabaseParam) (*sql.DB, error) {
	// hooks
	var driverHooks []hook.Hook

	// trace hook
	if p.Config.GetBool("modules.database.trace.enabled") {
		driverHooks = append(
			driverHooks,
			tracehook.NewTraceHook(
				tracehook.WithArguments(p.Config.GetBool("modules.database.trace.arguments")),
				tracehook.WithExcludedOperations(p.Config.GetStringSlice("modules.database.trace.exclude")...),
			),
		)
	}

	// log hook
	if p.Config.GetBool("modules.database.log.enabled") {
		driverHooks = append(
			driverHooks,
			loghook.NewLogHook(
				loghook.WithLevel(p.Config.GetString("modules.database.log.level")),
				loghook.WithArguments(p.Config.GetBool("modules.database.log.arguments")),
				loghook.WithExcludedOperations(p.Config.GetStringSlice("modules.database.log.exclude")...),
			),
		)
	}

	// driver registration
	driverName, err := RegisterDriver(p.Config.GetString("modules.database.driver"), driverHooks...)
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

type FxDatabaseMigratorParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
	Logger    *log.Logger
	Db        *sql.DB
}

func NewFxDatabaseMigrator(p FxDatabaseMigratorParam) (*migrate.Migrate, error) {
	var driver database.Driver
	var err error

	switch p.Config.GetString("modules.database.driver") {
	case "sqlite":
		driver, err = sqlite3migrate.WithInstance(p.Db, &sqlite3migrate.Config{})
	case "mysql":
		driver, err = mysqlmigrate.WithInstance(p.Db, &mysqlmigrate.Config{})
	case "postgres":
		driver, err = postgresmigrate.WithInstance(p.Db, &postgresmigrate.Config{})
	default:
		return nil, fmt.Errorf("invalid migration driver")
	}

	if err != nil {
		p.Logger.Error().Err(err).Msg("cannot build database migrations driver")

		return nil, err
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", p.Config.GetString("modules.database.migrations")),
		p.Config.GetString("modules.database.driver"),
		driver,
	)
	if err != nil {
		p.Logger.Error().Err(err).Msg("cannot build database migrations instance")

		return nil, err
	}

	fmt.Println("migrator up")

	migrator.Log = NewMigrationLogger(p.Logger)

	return migrator, nil
}

func RunFxDatabaseMigration(direction MigrationDirection, shutdown bool) fx.Option {
	return fx.Invoke(func(migrator *migrate.Migrate, logger *log.Logger, sd fx.Shutdowner) error {
		logger.Info().Msgf("starting database migrations (direction: %s)", direction)

		switch direction {
		case Up:
			err := migrator.Up()
			if err != nil {
				logger.Error().Err(err).Msg("cannot apply database migrations")

				return err
			}
		case Down:
			err := migrator.Down()
			if err != nil {
				logger.Error().Err(err).Msg("cannot apply databases migration")

				return err
			}
		default:
			err := fmt.Errorf("invalid migrations direction")
			logger.Error().Err(err).Msg("cannot apply database migrations")

			return err
		}

		logger.Info().Msgf("database migrations (direction: %s) applied with success", direction)

		if shutdown {
			return sd.Shutdown()
		}

		return nil
	})
}
