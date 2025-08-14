package main

import "snippetbox/internal/models"

//Define template data struct to hold our dynamic data
//Go has a limit of one per page, this allows us to do way more

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet //including a snippets field to hold a slice of snippets
}
