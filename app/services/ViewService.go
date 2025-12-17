package services

import (
	"bytes"
	"html/template"
	
	"log"
	"net/http"
	"path/filepath"
)

// View struct now holds the paths to be parsed on each render.
type View struct {
	paths []string
}

// NewViewService now only gathers the template paths and stores them.
func NewViewService(viewPaths ...string) *View {
	allPaths := []string{"app/views/layouts/", "app/views/components/"}
	allPaths = append(allPaths, viewPaths...)

	return &View{paths: allPaths}
}

// Render now parses templates on every call, enabling live reload.
func (v *View) Render(w http.ResponseWriter, name string, data interface{}) {
	var allFiles []string
	for _, path := range v.paths {
		files, err := filepath.Glob(filepath.Join(path, "*.html"))
		if err != nil {
			log.Printf("Error globbing templates: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		allFiles = append(allFiles, files...)
	}

	tmpl, err := template.ParseFiles(allFiles...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "app", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}