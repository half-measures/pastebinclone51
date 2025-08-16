package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	"snippetbox/internal/models"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" //special bit, when underscore we force import it
)

//Main is used for runtime config, dependencies for handlers and HTTP running

// Define our App struct to hold app wide dependencies,
// for now, just custom loggers
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	//remember ports 0-1023 are restricted
	addr := flag.String("addr", ":4000", "HTTP network address")
	// default value of 4000 set
	dsn := flag.String("dsn", "web:auxwork@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse() //Sanitizes the arg coming in just in case
	//we really really would want env vars but the drawback is no default setting out of the box
	//and no -help function

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// create logger for writing errs but we want stderr as dest
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close() //Always have, we want connection POOL to close before main func exits
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder() //init decoder instance to add to below dependencies

	// init a new instance of app struct for dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	//init a new server struct to use custom errorLog in problem event
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

	//Set Cache control header, if another Cache-Control header exists this will overwrite it

}

// OpenDB() function wraps sql.open and returns the sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn) //sql.open dosent create any connections, just inits a pool
	if err != nil {
		return nil, err

	} //ping to check if the connection is good
	if err = db.Ping(); err != nil {
		return nil, err

	}
	return db, nil
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

//Also have total control over which DB is used at runtime with -dsn cmd line flag

//DB -
//Generally its easy to swap out the DB, you can easily do so with GO
//still you have to remember to change the syntax
//Also note that GO is terrible with NULL, will throw a error as you cant convert it to string
//It does have a special type, but its easier to set constraints on DB columns to avoid null.
//Terraform TODO list

//1. Create mysql server with config, put into a env file and have go uptake that to get the secrets
//2.

//Routing notes
//Many to pick from, some have quirks
//best seem to be go-chi/chi, gorilla/mux, or julientschmidt/httprouter
//all three have good docs, tests, and work
//julien seems to be the most lightweight and fast, chi adds regexp patterns and groupings
//gorilla is most full featured yet slow.
