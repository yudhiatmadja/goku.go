package services

import (
    "bytes"
    "html/template"
    "io"
    "log"
    "net/http"
    "path/filepath"
)

type View struct {
    Template *template.Template
}

func NewViewService(viewPaths ...string) *View {
    // Tambahkan path default untuk layouts dan components
    paths := []string{"app/views/layouts/", "app/views/components/"}
    paths = append(paths, viewPaths...)

    var allFiles []string
    for _, path := range paths {
        files, err := filepath.Glob(filepath.Join(path, "*.html"))
        if err != nil {
            log.Fatalf("Failed to glob templates: %v", err)
        }
        allFiles = append(allFiles, files...)
    }

    tmpl := template.Must(template.ParseFiles(allFiles...))

    return &View{Template: tmpl}
}

func (v *View) Render(w io.Writer, name string, data interface{}) error {
    // Buat buffer untuk menangkap error sebelum menulis ke response
    buf := new(bytes.Buffer)
    err := v.Template.ExecuteTemplate(buf, name, data)
    if err != nil {
        return err
    }
    _, err = buf.WriteTo(w)
    return err
}