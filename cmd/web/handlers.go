package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// struct to represent form data for form fields
// struct must be exported and capitalized in order to be read by html/template package
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` //goes to Validators.go, embedding means this inherits all fields of the type Validator
}

// Our Handlers, it handels rendering stuff to user
// *http.request param is a pointer to a struct which holds info like http method and URL

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper.
	app.render(w, http.StatusOK, "home.tmpl", data)
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	// Use the new render helper.
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Display the form for creating a new snippet..."))
	data := app.newTemplateData(r)

	//init new createsnippetform instance pass to template
	data.Form = snippetCreateForm{
		Expires: 7,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// declare new empty instance of struct for decode method usage
	var form snippetCreateForm
	// call decode passing in current request and *pointer to createformstruct
	// this fills our struct will right values form html form, if err, do the thing
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//init a map to hold any validation errors from taking in the form fields

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank, fill it in now")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot more have than 100 chars long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank either, cmon dude")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "Field must equal 1, 7 or 365")
	// error check, dump any in plain http response and return
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
