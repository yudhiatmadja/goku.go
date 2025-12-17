package routes

import (
	"goku-framework/app/controllers"
	"goku-framework/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterWebRoutes(router *chi.Mux) {
	
	userController := &controllers.UserController{}
	homeController := controllers.NewHomeController()
	router.Get("/", homeController.Index)
	router.Get("/installation", homeController.Installation)
	router.Get("/docs", homeController.Docs)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)
		// r.Get("/profile", ProfileController.Show)
		// r.Get("/dashboard", DashboardController.Index)
	})

	router.Get("/users", userController.Index)
	// router.Post("/users", userController.Store)
}
