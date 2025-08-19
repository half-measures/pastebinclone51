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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
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
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "Field must equal 1, 7 or 365")
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
	// use put method to add string value and add key to session data
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Login Area funcs
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	//signing up a new user
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	//declare 0 value inst of the struct userSignupForm

	// Parse form data into said struct
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// validate using helper funcs
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	// if any errors, redisplay signup for wiht 422 statuscode
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}
	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}

		return
	}
	//otherwise add flash message to confirm signupworked
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in")
	//redirect page to login
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	// Hisplay html form for logging in user
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Auth and login the user
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// validation checks on form, check both email and pass provided and formatting of email
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}
	// check if cred are valid, if not add nonfielderr and goback2login
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// use renewtocken on current session to change sessID,
	// always generate a new session IF when auth state or priv changes for user
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// add id of current user to session so they are not logged in!
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	// redirect user to create a snippet
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// logout the user
	//userewnewtoken method on curr session id
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// remove auth userID from session data
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	// add flash message to session to tell user it worked
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	// redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
