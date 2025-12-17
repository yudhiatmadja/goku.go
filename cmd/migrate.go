package cmd

import (
	"fmt"
	"goku-framework/bootstrap"
	"goku-framework/database"

	"github.com/spf13/cobra"

	// Blank import to ensure init() functions in migrations are run
	// _ "goku-framework/database/migrations"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations based on migration files",
	Run: func(cmd *cobra.Command, args []string) {
		// bootstrap.NewApp() also initializes the database connection
		app := bootstrap.NewApp()

		if err := database.RunMigrations(app.DB); err != nil {
			fmt.Println("Migration failed:", err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}