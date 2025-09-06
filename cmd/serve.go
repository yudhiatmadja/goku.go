package cmd

import (
	"goku-framework/bootstrap"
	"goku-framework/routes"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Goku web server",
	Run: func(cmd *cobra.Command, args []string) {
		app := bootstrap.NewApp()
		// Daftarkan routes sebelum serve
		// routes.RegisterWebRoutes(app.Router)
		routes.RegisterWebRoutes(app.Router)
		app.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}