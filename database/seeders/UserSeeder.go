package seeders

import (
    "goku-framework/app/models"
    "gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
    users := []models.User{
        {Name: "Son Goku", Email: "goku@dbz.com"},
        {Name: "Vegeta", Email: "vegeta@dbz.com"},
    }
    
    // Hanya insert jika email belum ada
    for _, user := range users {
        db.FirstOrCreate(&user, models.User{Email: user.Email})
    }
    return nil
}