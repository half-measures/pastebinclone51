package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Routes method returns servemux with app routes
func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	//Handler function to wrap notFound helper, assign it as custom handler for 404 responses
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Update the pattern for the route for the static files.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// create new middleware chain for dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// And then create the routes using the appropriate methods, patterns and
	// handlers.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// Create the middleware chain as normal.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Wrap the router with the middleware and return it as normal.
	return standard.Then(router)
}

//Pas servemux as next param to secureheaders middleware
//secureheaders is a func and just returns http.Handler
