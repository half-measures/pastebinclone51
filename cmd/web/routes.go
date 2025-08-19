package main

import (
	"net/http"
	"snippetbox/ui"

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

	//

	// Update the pattern for the route for the static files.
	fileServer := http.FileServer(http.FS(ui.Files))
	//static files are in diff folder so we dont need to strip anything
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// create new middleware chain for dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate) //nosurf added for CSRF protection

	// And then create the routes using the appropriate methods, patterns and
	// handlers.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	//	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	//	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	//Five new routes w/ middleware for logging in users
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	//	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// protected auth only app routes w/ middleware chain
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create the middleware chain as normal.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Wrap the router with the middleware and return it as normal.
	return standard.Then(router)
}

//Pas servemux as next param to secureheaders middleware
//secureheaders is a func and just returns http.Handler
