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
	dri "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	driver := NewDriver(&dri.MySQLDriver{}, []hook.Hook{
		tracehook.NewTraceHook(),
		loghook.NewLogHook(),
	})

	driverName := fmt.Sprintf("yokai-%s", p.Config.GetString("modules.database.driver"))

	sql.Register(driverName, driver)

	db, err := sql.Open(driverName, p.Config.GetString("modules.database.dsn"))

	if err != nil {
		return nil, err
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return db.Close()
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
	driver, err := mysql.WithInstance(p.Db, &mysql.Config{})
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

	migrator.Log = NewMigrationLogger(p.Logger)

	return migrator, nil
}

func RunFxOrmDatabaseMigration(direction MigrationDirection) fx.Option {
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

		return sd.Shutdown()
	})
}
