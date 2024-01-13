package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gergab1129/bookings/internal/config"
	"github.com/gergab1129/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates string = "./templates"

// NewRenderer sets the config for template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template render templates using html
func Template(w http.ResponseWriter, tmpl string,
	td *models.TemplateData, r *http.Request) error {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// get template cache from app config

	// tc, err := CreateTemplateCache()

	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }

	// get requested template from cache
	t, ok := tc[tmpl]

	if !ok {
		fmt.Println("Template does not exists in cache")
		return errors.New("can't get template from cache")

	}

	td = AddDefaultData(td, r)

	buf := new(bytes.Buffer)

	err := t.Execute(buf, td)

	if err != nil {

		fmt.Println("Error: ", err)
		return err
	}
	// render the template

	_, err = buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	// parsedTemplate, err := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")

	// if err != nil {
	// 	fmt.Println("cannot parse file: ", err)
	// 	return
	// }

	// if err != nil {
	// 	fmt.Println("could not write template: ", err)
	// 	return
	// }

	return nil

}

func CreateTemplateCache() (map[string]*template.Template, error) {

	// templateCache := make(map[string]*template.Template)

	templateCache := map[string]*template.Template{}

	// get all of the files that start with .page.tmpk form ./templates

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		fmt.Println("Pattern not exists", err)
		return templateCache, err
	}

	// range for all the files matching *.pages.tmpl

	for _, page := range pages {
		fileName := filepath.Base(page)
		ts, err := template.New(fileName).ParseFiles(page)

		if err != nil {
			fmt.Println("Error: ", err)

		}

		layoutPath, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl",
			pathToTemplates))

		if err != nil {
			fmt.Println("Error: ", err)
		}

		if len(layoutPath) > 0 {

			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}

		templateCache[fileName] = ts
	}

	return templateCache, nil
}
