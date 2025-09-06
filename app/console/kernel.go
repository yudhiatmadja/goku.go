package console

import (
    "fmt"
    "github.com/robfig/cron/v3"
)

func RegisterScheduledTasks(c *cron.Cron) {
    // Contoh: jalankan setiap menit
    c.AddFunc("* * * * *", func() {
        fmt.Println("Running a scheduled task every minute!")
    })
}