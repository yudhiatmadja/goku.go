package routes

import (
	"goku-framework/app/controllers"
	"github.com/go-chi/chi/v5"
)

func RegisterApiRoutes(router *chi.Mux) {
	userController := &controllers.UserController{}

	router.Route("/api", func(r chi.Router) {
		r.Get("/users", userController.Index)
		r.Post("/users", userController.Store)
	})
}
