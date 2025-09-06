package bootstrap

import "github.com/hibiken/asynq"

// Konfigurasi Redis dari config
var RedisOpt = asynq.RedisClientOpt{Addr: "localhost:6379"}