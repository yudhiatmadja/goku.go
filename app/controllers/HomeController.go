package controllers

import (
    "goku-framework/app/services"
    "net/http"
)

type HomeController struct {
    View *services.View
}

func NewHomeController() *HomeController {
    // Path ke folder pages
    viewService := services.NewViewService("app/views/pages/")
    return &HomeController{View: viewService}
}

func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "Title": "Homepage",
        "Name":  "Goku User",
    }
    c.View.Render(w, "home.html", data)
}