package cmd

import (
    "goku-framework/app/console"
    "github.com/robfig/cron/v3"
    "github.com/spf13/cobra"
)

var scheduleRunCmd = &cobra.Command{
    Use:   "schedule:run",
    Short: "Run the scheduled tasks",
    Run: func(cmd *cobra.Command, args []string) {
        c := cron.New()
        console.RegisterScheduledTasks(c)
        c.Start()
        select {} // Block forever
    },
}

func init() {
    rootCmd.AddCommand(scheduleRunCmd)
}