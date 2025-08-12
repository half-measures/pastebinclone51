package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

//Main is used for runtime config, dependencies for handlers and HTTP running

// Define our App struct to hold app wide dependencies,
// for now, just custom loggers
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	//remember ports 0-1023 are restricted
	addr := flag.String("addr", ":4000", "HTTP network address")
	// default value of 4000 set
	flag.Parse() //Sanitizes the arg coming in just in case
	//we really really would want env vars but the drawback is no default setting out of the box
	//and no -help function

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// create logger for writing errs but we want stderr as dest
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// init a new instance of app struct for dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	//init a new server struct to use custom errorLog in problem event
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

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
