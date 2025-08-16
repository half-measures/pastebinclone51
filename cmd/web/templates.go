package main

import (
	"path/filepath"
	"snippetbox/internal/models"
	"text/template"
	"time"
)

//Define template data struct to hold our dynamic data
//Go has a limit of one per page, this allows us to do way more

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet //including a snippets field to hold a slice of snippets
	Form        any               //Used to pass validation errors back to template when re-display form so users dont have to enter it again
}

// Formating a nicer string for time
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// init a funcmap object and store it in global var
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the  page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal...
		cache[name] = ts
	}

	return cache, nil
}
