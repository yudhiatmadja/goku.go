package routes

import (
	"goku-framework/app/controllers"
	"goku-framework/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterWebRoutes(router *chi.Mux) {
	
	homeController := controllers.NewHomeController()
	router.Get("/", homeController.Index)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)

	})


}
