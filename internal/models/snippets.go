package models

import (
	"database/sql"
	"errors"
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
	//SQL
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10"

	//Use Query method to exec, returns more than one row tho!
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	//Ensure prop closed, after connection
	//Always have, if not, connection is open and many things start to go wrong as you scale

	defer rows.Close()
	snippets := []*Snippet{}
	//init empty slice to hold struct

	//loop thru the resultset
	for rows.Next() {
		s := &Snippet{}
		//we use rows.scan to copy values from each field in the row to a new object we created.
		//new objects must be pointers
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		//Append to slice of snippets
		snippets = append(snippets, s)
	}

	//When loop finished, call rows.err to get any err if any
	//Never assume success with Databases, ensure it
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
	//If ok, return snippets slice
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	//SQL
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"
	//use of query row instead on conn pool as we only want a single row result
	row := m.DB.QueryRow(stmt, id)
	s := &Snippet{}
	//Pointer to a new zeroed snippet struct

	//Use row.Scan() to copy values from each field in row to snip struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		//If Query returns no rows, scan returns a ErrNowRows error.
		//We should use errors.IS function to check for that err
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	//if everything is OK, then return our snip object
	return s, nil
}

//Were adding this new snippet struct to represent data for snippet along with our
//snippet model type - Need to add to main.go and inject it as a dependecies
//cuz of how this is set, db logic is not around our handlers whihc means
//those are still just for HTTP stuff to make it easy for unit tests
//We also have
