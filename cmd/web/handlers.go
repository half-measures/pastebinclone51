package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

// Our Handlers, it handels rendering stuff to user
// *http.request param is a pointer to a struct which holds info like http method and URL
func home(w http.ResponseWriter, r *http.Request) {
	// add 404
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
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
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	//We use Execute() on template set to write template content as reponse body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}
func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	//use the printf to interpolate the id value with our response
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	//fprintf is a type interface and the http.responsewriter object satifys the req of  interface needed as it has a w.write() method
	//For now, when we see a io.Writer Param, its ok to pass your responseWriterobject.
	//Whatever being written will be sent as the body of the HTTP response.
}

func snippetCreate(w http.ResponseWriter, r *http.Request) { //
	//use r.Method to check if its post or not, POST causes a change to server so should be only way to do that
	if r.Method != http.MethodPost {
		//if its not, we use w.WriteHeader() to send 405 status (Method not allowed)
		//This is good security to avoid being hit with errent hits, and in the future this will be a DB call
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet zone..."))
}
