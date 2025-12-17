package cmd

import (
    "goku-framework/app/jobs"
    "goku-framework/bootstrap"
    "log"
    "github.com/hibiken/asynq"
    "github.com/spf13/cobra"
)

var queueWorkCmd = &cobra.Command{
    Use:   "queue:work",
    Short: "Start processing jobs on the queue",
    Run: func(cmd *cobra.Command, args []string) {
        srv := asynq.NewServer(bootstrap.RedisOpt, asynq.Config{Concurrency: 10})
        
        mux := asynq.NewServeMux()
        mux.HandleFunc(jobs.TypeWelcomeEmail, jobs.HandleWelcomeEmailTask)

        if err := srv.Run(mux); err != nil {
            log.Fatalf("could not run server: %v", err)
        }
    },
}

func init() {
    RootCmd.AddCommand(queueWorkCmd)
}