package bootstrap

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"goku-framework/database"
)

var Router *chi.Mux

type Application struct {
	Config *viper.Viper
	DB     *gorm.DB
	Router *chi.Mux
}

func NewApp() *Application {
	app := &Application{}
	app.initConfig()
	app.initDB()
	app.initRouter()
	return app
}

func InitApp() {
	Router = chi.NewRouter()

	// contoh route
	Router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("List users"))
	})

	Router.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Create user"))
	})
}

func (app *Application) initConfig() {
	app.Config = viper.New()
	app.Config.AddConfigPath("./config") // Path ke folder config
	app.Config.SetConfigName("app")      // Nama file tanpa ekstensi
	app.Config.SetConfigType("yaml")

	if err := app.Config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func (app *Application) initDB() {
	dsn := app.Config.GetString("database.dsn")
	dbName := app.Config.GetString("database.name")

	if dsn == "" || dbName == "" {
		log.Fatalf("Database configuration is missing in config/app.yaml")
	}

	// Create database if it doesn't exist
	dsnWithoutDB := strings.Replace(dsn, "/"+dbName, "/", 1)
	sqlDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		log.Fatalf("Failed to open connection to MySQL server: %v", err)
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Connect to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	app.DB = db
	database.DB = db
}

func (app *Application) initRouter() {
	app.Router = chi.NewRouter()
	// Di sini nanti kita akan mendaftarkan routes dari routes/web.go dan routes/api.go
}

func (app *Application) Serve() {
	port := app.Config.GetString("server.port")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Printf("Goku server running on port %s\n", port)
	http.ListenAndServe(":"+port, app.Router)
}

