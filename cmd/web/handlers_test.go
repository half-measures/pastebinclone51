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

// Some of this func now in testutils_test.go
// adding more tests to see req body contains some content, but not exaxtly equal to it
func TestSnippetView(t *testing.T) {
	//create new inst of app struct with mock dependents
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	//set up table driven tests to check responses
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
