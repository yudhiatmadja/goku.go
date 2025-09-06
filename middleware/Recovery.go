package middleware

import (
    "fmt"
    "net/http"
    "github.com/sirupsen/logrus"
)

func Recoverer(logger *logrus.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    logger.WithFields(logrus.Fields{
                        "error": err,
                        "stack": string(debug.Stack()), // perlu import "runtime/debug"
                    }).Error("Panic recovered")

                    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}