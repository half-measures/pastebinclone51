package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"strconv"
	"text/template"
)

// Our Handlers, it handels rendering stuff to user
// *http.request param is a pointer to a struct which holds info like http method and URL
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// add 404
	if r.URL.Path != "/" {
		app.notFound(w) //Use our helpers
		return
	}
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	//init a slice to combine our two files
	//note that base template must come first in slice
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	//Use template.ParseFiles() to read our html template file into a set
	//If err, we log and send a 500 Internal Server error response
	ts, err := template.ParseFiles(files...) //Relative to root of project instead of Hardcoded
	if err != nil {
		app.serverError(w, err) //Use our helper
		return
	}

	//We use Execute() on template set to write template content as reponse body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err) //use the serverError() helper
	}

}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) //see helper.go
		return
	}
	//use the snippetmodel object get method to get data for a record based on its ID
	//if none found, retrn 404 not found
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	//write snip data back as plaintext HTTP
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
