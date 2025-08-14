package main

import (
	"path/filepath"
	"snippetbox/internal/models"
	"text/template"
)

//Define template data struct to hold our dynamic data
//Go has a limit of one per page, this allows us to do way more

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet //including a snippets field to hold a slice of snippets
}

func newTemplateCache() (map[string]*template.Template, error) {
	// init new map to act as a cache
	cache := map[string]*template.Template{}

	// use the filepath.glb to get a slice of all filepaths that match the pattern
	// gives us a slice of all filepaths that are in our apps templates folder
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	// loop thru page filepaths one by one
	for _, page := range pages {
		//extract the file name like home.tmpl from filepath
		name := filepath.Base(page)
		//create slice w/ filepaths for base templates. include partials
		files := []string{
			"./ui/html/base.tmpl",
			".ui/html/partials/nav.tmpl",
			page,
		}
		//parse files into templateset
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		//add template set to map
		cache[name] = ts
	}
	// return map
	return cache, nil
}
