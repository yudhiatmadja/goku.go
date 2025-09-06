package bootstrap

import (
    "github.com/sirupsen/logrus"
    "gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *logrus.Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})

    // Konfigurasi rotasi log harian
    logFile := &lumberjack.Logger{
        Filename:   "storage/logs/goku.log",
        MaxSize:    5, // megabytes
        MaxBackups: 10,
        MaxAge:     28, // days
        Compress:   true,
    }
    
    logger.SetOutput(logFile)
    return logger
}