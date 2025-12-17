package cmd

import (
	"goku-framework/bootstrap"
	"goku-framework/middleware"
	"goku-framework/routes"
	"goku-framework/util"
	"log"

	"github.com/spf13/cobra"
)

var watch bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Goku web server",
	Run: func(cmd *cobra.Command, args []string) {
		app := bootstrap.NewApp()

		if watch {
			log.Println("Goku server starting in development mode with live-reload enabled.")

			// 1. Create and run the WebSocket hub
			hub := util.NewHub()
			go hub.Run()

			// 2. Start the file watcher in a goroutine
			go util.WatchFiles(hub)

			// 3. Apply the live-reload middleware FIRST
			app.Router.Use(middleware.LiveReload)

			// 4. Register the WebSocket endpoint
			app.Router.Get("/goku-ws", hub.ServeWs)
		}

		routes.RegisterWebRoutes(app.Router)
		routes.RegisterApiRoutes(app.Router)
		app.Serve()
	},
}

func init() {
	serveCmd.Flags().BoolVar(&watch, "watch", false, "Enable watch mode for live reloading")
	RootCmd.AddCommand(serveCmd)
}