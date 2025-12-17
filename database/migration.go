package database

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"sort"
)

// Migration defines the structure for a single database migration.
type Migration struct {
	Timestamp string
	Name      string
	Up        func(*gorm.DB) error
	Down      func(*gorm.DB) error
}

// SchemaMigration is the model for the migrations table in the database.
type SchemaMigration struct {
	Timestamp string `gorm:"primaryKey"`
}

var migrations []Migration

// AddMigration is called by the init() function of each migration file.
func AddMigration(m Migration) {
	migrations = append(migrations, m)
}

// RunMigrations executes all pending migrations.
func RunMigrations(db *gorm.DB) error {
	// 1. Ensure the schema_migrations table exists.
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return fmt.Errorf("failed to auto-migrate schema_migrations table: %w", err)
	}

	// 2. Get all migrations that have already been run.
	var ranMigrations []SchemaMigration
	if err := db.Find(&ranMigrations).Error; err != nil {
		return fmt.Errorf("failed to get ran migrations: %w", err)
	}

	ranMap := make(map[string]bool)
	for _, m := range ranMigrations {
		ranMap[m.Timestamp] = true
	}

	// 3. Sort migrations by timestamp.
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Timestamp < migrations[j].Timestamp
	})

	// 4. Run pending migrations.
	log.Println("Running pending migrations...")
	for _, m := range migrations {
		if !ranMap[m.Timestamp] {
			log.Printf("Applying migration: %s_%s\n", m.Timestamp, m.Name)
			if err := m.Up(db); err != nil {
				return fmt.Errorf("failed to run migration '%s': %w", m.Name, err)
			}

			// 5. Record the migration as run.
			newMigration := SchemaMigration{Timestamp: m.Timestamp}
			if err := db.Create(&newMigration).Error; err != nil {
				return fmt.Errorf("failed to record migration '%s' as run: %w", m.Name, err)
			}
			log.Printf("Successfully applied migration: %s_%s\n", m.Timestamp, m.Name)
		}
	}

	log.Println("All migrations have been run successfully.")
	return nil
}
