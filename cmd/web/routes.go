package main

import "net/http"

// Routes method returns servemux with app routes
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.logRequest(secureHeaders(mux))
}

//Pas servemux as next param to secureheaders middleware
//secureheaders is a func and just returns http.Handler
