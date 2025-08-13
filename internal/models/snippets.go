package models

import (
	"database/sql"
	"time"
)

// define a snippet type to hold data for indiv snippet.
// Fields must correspond to fields in our SQL snips
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
} //Defines snip model to wrap a sql connection pool

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	//pgsql uses $N, msql uses ?
	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// title, content and expiry values for the placeholder parameters. This
	// method returns a sql.Result type, which contains some basic
	// information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// This returns 10 most recent snips
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}

//Were adding this new snippet struct to represent data for snippet along with our
//snippet model type - Need to add to main.go and inject it as a dependecies
//cuz of how this is set, db logic is not around our handlers whihc means
//those are still just for HTTP stuff to make it easy for unit tests
//We also have
