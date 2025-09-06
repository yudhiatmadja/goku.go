package bootstrap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	// Contoh koneksi DB dari config (perlu setup database.yaml)
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// Untuk sekarang, kita hardcode dulu
	dsn := "root:@tcp(127.0.0.1:3306)/goku_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	app.DB = db
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

