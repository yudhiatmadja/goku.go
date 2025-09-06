package controllers

import (
	"encoding/json"
	"net/http"
	"goku-framework/app/http/requests"
    "goku-framework/app/models"
    "gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// Index - Menampilkan daftar user
func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	// Logika untuk mengambil data user dari database
	// Untuk sekarang, kita return data dummy
	users := []map[string]string{
		{"id": "1", "name": "Son Goku"},
		{"id": "2", "name": "Vegeta"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *UserController) Store(w http.ResponseWriter, r *http.Request) {
    var req requests.UserCreateRequest
    if err := requests.Validate(r, &req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user := models.User{Name: req.Name, Email: req.Email}
    result := c.DB.Create(&user)

    if result.Error != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("User created successfully"))
}