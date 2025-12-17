package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var migrationTemplate = `package migrations

import (
	"goku-framework/database"
	"gorm.io/gorm"
	"log"
	// "goku-framework/app/models" // Example: import models
)

func init() {
	database.AddMigration(database.Migration{
		Timestamp: "{{.Timestamp}}",
		Name:      "{{.MigrationName}}",
		Up:        up{{.StructName}},
		Down:      down{{.StructName}},
	})
}

func up{{.StructName}}(db *gorm.DB) error {
	log.Println("Migrating up: {{.MigrationName}}")
	// Example: return db.AutoMigrate(&models.User{})
	return nil
}

func down{{.StructName}}(db *gorm.DB) error {
	log.Println("Migrating down: {{.MigrationName}}")
	// Example: return db.Migrator().DropTable(&models.User{})
	return nil
}
`

var MakeMigrationCmd = &cobra.Command{
	Use:   "make:migration [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")
		fileName := fmt.Sprintf("%s_%s.go", timestamp, name)
		structName := strcase.ToCamel(name)

		// Create file
		filePath := filepath.Join("database", "migrations", fileName)
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating migration file: %v\n", err)
			return
		}
		defer file.Close()

		// Parse template
		tmpl, err := template.New("migration").Parse(migrationTemplate)
		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			return
		}

		// Prepare data for template
		data := struct {
			Timestamp     string
			MigrationName string
			StructName    string
		}{
			Timestamp:     timestamp,
			MigrationName: name,
			StructName:    structName,
		}

		if err := tmpl.Execute(file, data); err != nil {
			fmt.Printf("Error executing template: %v\n", err)
			return
		}

		fmt.Printf("Migration [%s] created successfully.\n", filePath)
	},
}

func init() {
	RootCmd.AddCommand(MakeMigrationCmd)
}
