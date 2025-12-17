package cmd

import (
    "fmt"
    "goku-framework/bootstrap"
    "goku-framework/database/seeders"
    "github.com/spf13/cobra"
)

var dbSeedCmd = &cobra.Command{
    Use:   "db:seed",
    Short: "Seed the database with records",
    Run: func(cmd *cobra.Command, args []string) {
        app := bootstrap.NewApp()
        fmt.Println("Seeding database...")

        // Panggil semua seeder di sini
        if err := seeders.SeedUsers(app.DB); err != nil {
            fmt.Println("Failed to seed users:", err)
        }

        fmt.Println("Database seeding completed successfully.")
    },
}

func init() {
    RootCmd.AddCommand(dbSeedCmd)
}