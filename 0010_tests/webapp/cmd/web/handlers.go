package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

var pathToTemplate = "./templates/"

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	fmt.Println("->>> ", path.Join(pathToTemplate), path.Join(pathToTemplate)+t)
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplate, t))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}
	return parsedTemplate.Execute(w, data)
}
