package cmd

import (
	"fmt"
	"goku-framework/bootstrap"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

// RootCmd is the root command for the application.
var RootCmd = &cobra.Command{
	Use:   "goku",
	Short: "Goku is a Laravel-like Go web framework",
	Long:  `A fast, simple, and productive web framework for Go developers.`,
}

var routeListCmd = &cobra.Command{
	Use:   "route:list",
	Short: "Menampilkan daftar route",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.InitApp()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "METHOD\tPATH")

		chi.Walk(bootstrap.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			fmt.Fprintf(w, "%s\t%s\n", method, route)
			return nil
		})

		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(routeListCmd)
	// Other commands from the 'cmd' package are added in their own files
}

// Execute runs the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}