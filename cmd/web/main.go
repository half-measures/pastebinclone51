package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	//remember ports 0-1023 are restricted
	addr := flag.String("addr", ":4000", "HTTP network address")
	// default value of 4000 set
	flag.Parse() //Sanitizes the arg coming in just in case

	// http.NewServeMux() func to init a new servemux then register it
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home) //note that / servers all, so any path goes to a single place with that right now
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/Create", snippetCreate) //Note Capitilzation
	//use http.ListenAndServe to start new web server, passing twoparams
	//like our port number

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)

	//Set Cache control header, if another Cache-Control header exists this will overwrite it

}

//Web App basics include a handler - its a bit like a controller and do app logic
//write HTTP responses and headers
//2nd thing we need is a router to store URL maps like serveux

//To be consise, when this server above gets a new http request, it calles the servemux
//ServeHTTP() method which is abstracted away from us
//It finds the right handler based on request URL path and calls
//that handlers ServeHTTP() method
//In a way, this is all a chain of ServeHTTP() methods being called by one another

//Also, all HTTP connections are served via there own goroutine
//This makes it very fast but we need to be mindful of race conditions in the future.
