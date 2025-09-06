package cmd

import (
    "fmt"
    "goku-framework/app/models"
    "goku-framework/bootstrap"
    "github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Run database migrations",
    Run: func(cmd *cobra.Command, args []string) {
        app := bootstrap.NewApp() // Hanya untuk koneksi DB
        fmt.Println("Running migrations...")
        
        // Daftarkan semua model di sini
        err := app.DB.AutoMigrate(&models.User{})
        if err != nil {
            fmt.Println("Migration failed:", err)
            return
        }

        fmt.Println("Migration completed successfully.")
    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}