package routes

import (
	"goku-framework/app/controllers"

	"github.com/go-chi/chi/v5"
)

func RegisterWebRoutes(router *chi.Mux) {
	// Buat instance dari controller
	userController := &controllers.UserController{}
	homeController := controllers.NewHomeController()
	router.Get("/", homeController.Index)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Goku Framework!"))
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/profile", ProfileController.Show)
		r.Get("/dashboard", DashboardController.Index)
	})

	router.Get("/users", userController.Index)
	// router.Post("/users", userController.Store)
}
