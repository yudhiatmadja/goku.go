package migrations

import (
	"goku-framework/database"
	"gorm.io/gorm"
	"log"
	// "goku-framework/app/models" // Example: import models
)

func init() {
	database.AddMigration(database.Migration{
		Timestamp: "20251217140129",
		Name:      "create_users_table",
		Up:        upCreateUsersTable,
		Down:      downCreateUsersTable,
	})
}

func upCreateUsersTable(db *gorm.DB) error {
	log.Println("Migrating up: create_users_table")
	// Example: return db.AutoMigrate(&models.User{})
	return nil
}

func downCreateUsersTable(db *gorm.DB) error {
	log.Println("Migrating down: create_users_table")
	// Example: return db.Migrator().DropTable(&models.User{})
	return nil
}
