package main

//test our http runnin's
import (
	"net/http"
	"snippetbox/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	//new init of app struct for mock logging for now
	app := newTestApplication(t)
	//NewTLSserver for new test server, using app.routes method, tests
	//random part and shutsdown after test is done
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	//network address that testserv is listening on is in ts.url field,
	//use get .ping to request against test server with http.Response returning
	code, _, body := ts.get(t, "/ping")

	//check value of status code
	assert.Equal(t, code, http.StatusOK)

	assert.Equal(t, body, "OK")

}

//Some of this func now in testutils_test.go
