package cmd

import (
	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase"
	"github.com/ankorstore/yokai/fxcore"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down]",
	Short: "Run application database migrations",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fxcore.
			NewBootstrapper().
			WithContext(cmd.Context()).
			WithOptions(
				fx.NopLogger,
				// modules
				fxdatabase.FxDatabaseModule,
				// migration
				fxdatabase.RunFxOrmDatabaseMigration(fxdatabase.FetchMigrationDirection(args[0])),
			).
			RunApp()
	},
}
