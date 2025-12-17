# Goku Framework

<p align="center">
  <img src="https://raw.githubusercontent.com/gtkm/goku-framework/main/storage/public/goku-logo.png" width="200" alt="Goku Framework Logo">
</p>

<p align="center">
  A modern Go web framework built for developers who value simplicity, performance, and elegant code architecture.
</p>

<p align="center">
  <a href="https://github.com/gtkm/goku-framework/actions"><img src="https://github.com/gtkm/goku-framework/actions/workflows/go.yml/badge.svg" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/gtkm/goku-framework"><img src="https://goreportcard.com/badge/github.com/gtkm/goku-framework" alt="Go Report Card"></a>
  <a href="https://github.com/gtkm/goku-framework/blob/main/LICENSE"><img src="https://img.shields.io/github/license/gtkm/goku-framework.svg" alt="License"></a>
</p>

Goku is a web framework for the Go programming language that is inspired by the elegance and developer-friendliness of frameworks like Laravel. It provides a solid foundation for building robust and performant web applications with a clean and maintainable structure.

## Features

- **Fast Routing**: Built on top of the lightweight and high-performance `chi` router.
- **Powerful ORM**: Seamless integration with `GORM` for a developer-friendly database interaction.
- **Built-in CLI**: An `Artisan`-like command-line tool for generating boilerplate and managing your application.
- **Live Reload**: A `--watch` mode to automatically reload the server on file changes.
- **Organized Structure**: A logical and predictable directory structure.
- **Easy Configuration**: Effortless configuration management with `Viper`.
- **Session Management**: Built-in support for sessions using `gorilla/sessions`.
- **Job Queuing**: A robust job queuing system powered by `asynq` and Redis.
- **Task Scheduling**: A simple way to schedule recurring tasks using a cron-based scheduler.

---

## Getting Started

### Installation

To start a new project with Goku Framework, clone this repository:

```bash
git clone https://github.com/gtkm/goku-framework.git my-goku-app
cd my-goku-app
go mod tidy
```

### Configuration

1.  Copy the example configuration file.
    ```bash
    cp config/app.example.yaml config/app.yaml
    ```
2.  Edit `config/app.yaml` to match your local environment, especially the database connection details.

    ```yaml
    app:
      name: Goku
      port: 8080
      env: development

    database:
      driver: "mysql"
      dsn: "user:password@tcp(127.0.0.1:3306)/" # Note: No database name here
      name: "goku_db" # The database name to be created or used
    
    session:
      secret: "change-this-to-a-very-secret-key"

    redis:
      addr: "localhost:6379"
    ```

### Running the Application

To run the development server:

```bash
go run main.go serve
```

For development with live-reload, use the `--watch` flag:

```bash
go run main.go serve --watch
```

The server will be available at `http://localhost:8080`.

---

## Core Concepts

### Project Structure

The Goku directory structure is designed to be intuitive and scalable.

```
.
├── app
│   ├── console       # Scheduled tasks (cron jobs)
│   ├── controllers   # HTTP controllers
│   ├── http          # Form Requests and validation
│   ├── jobs          # Queueable jobs
│   ├── models        # GORM models
│   └── views         # HTML templates
├── bootstrap         # Application initialization (DB, session, etc.)
├── cmd               # CLI (Cobra) commands
├── config            # Configuration files
├── database
│   ├── migrations    # Database migration files
│   └── seeders       # Database seeder files
├── middleware        # HTTP middleware
├── routes            # Route definitions
├── storage           # File storage (public, framework, app)
└── util              # Common utilities
```

### Routing

Routes are defined in `routes/web.go`. Goku uses `chi` for routing.

**Route Parameters:**
```go
r.Get("/users/{id}", userController.Show)
```

**Route Groups & Middleware:**
```go
router.Group(func(r chi.Router) {
    r.Use(middleware.Authenticate) // Protects all routes in this group
    r.Get("/profile", ProfileController.Show)
    r.Get("/settings", SettingsController.Index)
})
```

### Static File Serving

Goku does not serve static files by default. You need to enable it manually. This gives you control over which directories are publicly accessible.

To serve files from the `storage/public` directory, add the following to `cmd/serve.go` before `app.Serve()`:

```go
// In cmd/serve.go, inside the Run function for serveCmd
// ...
routes.RegisterWebRoutes(app.Router)

// Add this block for static file serving
fs := http.FileServer(http.Dir("./storage/public"))
app.Router.Handle("/static/*", http.StripPrefix("/static/", fs))

app.Serve()
// ...
```
Now, a file at `storage/public/css/style.css` will be accessible at `http://localhost:8080/static/css/style.css`.

### Session Management

Goku uses `gorilla/sessions` for session management. The session store is initialized in `bootstrap/session.go`.

You can access the session store from any controller:

**Setting a Session Value:**
```go
session, _ := bootstrap.Store.Get(r, "goku-session")
session.Values["user_id"] = 123
err := session.Save(r, w)
```

**Getting a Session Value:**
```go
session, _ := bootstrap.Store.Get(r, "goku-session")
userID := session.Values["user_id"]
```

### Authentication

The framework includes a basic `Authenticate` middleware that checks for a session value.

**How It Works:**
The middleware in `middleware/Auth.go` checks if `session.Values["authenticated"]` is `true`. If not, it redirects the user to `/login`.

**Login Controller Example:**
```go
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
    // ... (validate user credentials)

    // If valid:
    session, _ := bootstrap.Store.Get(r, "goku-session")
    session.Values["authenticated"] = true
    session.Values["user_id"] = foundUser.ID
    session.Save(r, w)

    http.Redirect(w, r, "/dashboard", http.StatusFound)
}
```

**Logout Controller Example:**
```go
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := bootstrap.Store.Get(r, "goku-session")
    session.Values["authenticated"] = false
    session.Options.MaxAge = -1 // Deletes the cookie
    session.Save(r, w)
    
    http.Redirect(w, r, "/login", http.StatusFound)
}
```

### Request Validation

Goku uses `go-playground/validator` for request validation. Create request structs in `app/http/requests`.

**Example Request Struct:**
```go
// app/http/requests/UserCreateRequest.go
package requests

type UserCreateRequest struct {
    Name  string `json:"name" validate:"required,min=3"
    Email string `json:"email" validate:"required,email"
}
```

The `Validate` helper in `app/http/requests/validator.go` handles decoding and validation.

---

## Database

### Migrations

Create a migration with the `make:migration` command.
```bash
go run main.go make:migration create_users_table
```
Edit the generated file in `database/migrations` to define your schema. Then, run the migrations:
```bash
go run main.go migrate
```

### Seeding

Seeders in `database/seeders` populate your database.
```bash
go run main.go db:seed
```

---

## Advanced Topics

### Job Queuing

Goku uses `asynq` with a Redis backend for background jobs.

**1. Define a Job:**
Create a job definition in `app/jobs`. A job needs a `Type` and a `Payload`.

```go
// app/jobs/SendWelcomeEmail.go
package jobs

const TypeWelcomeEmail = "email:welcome"

type WelcomeEmailPayload struct {
    UserID int
}

func NewWelcomeEmailTask(userID int) (*asynq.Task, error) {
    // ... (implementation)
}
```

**2. Create a Handler:**
The handler processes the job.

```go
// app/jobs/SendWelcomeEmail.go
func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
    // ... (implementation)
}
```

**3. Register the Handler:**
In `cmd/queue.go`, register your handler.
```go
// cmd/queue.go
mux := asynq.NewServeMux()
mux.HandleFunc(jobs.TypeWelcomeEmail, jobs.HandleWelcomeEmailTask)
```

**4. Dispatch the Job:**
From your application (e.g., in a controller), you can dispatch the job.

```go
client := asynq.NewClient(bootstrap.RedisOpt)
task, err := jobs.NewWelcomeEmailTask(newUser.ID)
if err != nil {
    // Handle error
}
info, err := client.Enqueue(task)
if err != nil {
    // Handle error
}
```

**5. Run the Worker:**
Start a queue worker to process the jobs.
```bash
go run main.go queue:work
```

### Task Scheduling

You can schedule recurring tasks in `app/console/kernel.go` using a cron expression.

```go
// app/console/kernel.go
package console

import (
    "fmt"
    "github.com/robfig/cron/v3"
)

func RegisterScheduledTasks(c *cron.Cron) {
    // Run a task every minute
    c.AddFunc("* * * * *", func() {
        fmt.Println("Running a scheduled task!")
    })
    
    // Run a task every day at midnight
    c.AddFunc("@daily", func() { /* ... */ })
}
```
To run the scheduler:
```bash
go run main.go schedule:run
```

### Custom CLI Commands

You can create your own CLI commands, similar to `serve` or `migrate`.

1.  Create a new file in the `cmd` directory (e.g., `cmd/my_command.go`).
2.  Define a `cobra.Command` and add it to the `RootCmd`.

**Example `cmd/greet.go`:**
```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var greetCmd = &cobra.Command{
    Use:   "greet [name]",
    Short: "Greets a person",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Hello, %s!\n", args[0])
    },
}

func init() {
    RootCmd.AddCommand(greetCmd)
}
```

Now you can run your command:
```bash
go run main.go greet World
```

---

## Contributing

Contributions are welcome! Please fork the repository, create a feature branch, and submit a pull request.

## License

Goku Framework is open-source software licensed under the [MIT License](https://github.com/gtkm/goku-framework/blob/main/LICENSE).
