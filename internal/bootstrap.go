package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxsql"
	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxsqlc"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

func init() {
	RootDir = fxcore.RootDir(1)
}

// RootDir is the application root directory.
var RootDir string

// Bootstrapper can be used to load modules, options, dependencies, routing and bootstraps the application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxhttpserver.FxHttpServerModule,
	fxsql.FxSQLModule,
	fxsqlc.FxSqlcModule,
	// dependencies registration
	Register(),
	// routing registration
	Router(),
)

// Run starts the application, with a provided [context.Context].
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp()
}

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	tb.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", RootDir))
	tb.Setenv("MODULES_DATABASE_MIGRATIONS", fmt.Sprintf("%s/db/migrations", RootDir))

	Bootstrapper.RunTestApp(tb, fx.Options(options...))
}
