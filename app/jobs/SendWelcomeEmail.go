package jobs

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/hibiken/asynq"
)

const TypeWelcomeEmail = "email:welcome"

type WelcomeEmailPayload struct {
    UserID int
}

func NewWelcomeEmailTask(userID int) (*asynq.Task, error) {
    payload, err := json.Marshal(WelcomeEmailPayload{UserID: userID})
    if err != nil {
        return nil, err
    }
    return asynq.NewTask(TypeWelcomeEmail, payload), nil
}

func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
    var p WelcomeEmailPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("failed to unmarshal payload: %w", err)
    }
    fmt.Printf("Sending welcome email to user %d\n", p.UserID)
    // Logika pengiriman email di sini...
    return nil
}