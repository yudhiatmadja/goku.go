package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var controllerTemplate = `package controllers

import (
	"net/http"
)

type {{.ControllerName}} struct {
}

// Index method
func (c *{{.ControllerName}}) Index(w http.ResponseWriter, r *http.Request) {
	// Your logic here
}

// Store method
func (c *{{.ControllerName}}) Store(w http.ResponseWriter, r *http.Request) {
	// Your logic here
}
`

var MakeControllerCmd = &cobra.Command{
	Use:   "make:controller [name]",
	Short: "Create a new controller file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if !strings.HasSuffix(name, "Controller") {
			name += "Controller"
		}

		// Buat file
		filePath := filepath.Join("app", "controllers", name+".go")
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating controller file: %v\n", err)
			return
		}
		defer file.Close()

		// Parse template
		tmpl, err := template.New("controller").Parse(controllerTemplate)
		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			return
		}

		// Tulis ke file
		data := struct {
			ControllerName string
		}{
			ControllerName: name,
		}
		if err := tmpl.Execute(file, data); err != nil {
			fmt.Printf("Error executing template: %v\n", err)
			return
		}

		fmt.Printf("Controller [%s] created successfully.\n", filePath)
	},
}

func init() {
	RootCmd.AddCommand(MakeControllerCmd)
}
