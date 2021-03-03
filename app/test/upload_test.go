package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joerx/lolcatz-backend/test/dbtest"
)

func Test_listUploads(t *testing.T) {
	// truncate database
	if err := dbtest.Truncate(dbc); err != nil {
		t.Fatal(err)
	}

	// seed uploads
	uploads, err := dbtest.SeedUploads(dbc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/list"), nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if e, a := http.StatusOK, w.Result().StatusCode; e != a {
		t.Errorf("Expected statuscode %d but got %d", e, a)
	}

	body, _ := ioutil.ReadAll(w.Body)
	items := make([]map[string]interface{}, 0)

	if err := json.Unmarshal(body, &items); err != nil {
		t.Fatal(err)
	}

	if e, a := len(uploads), len(items); e != a {
		t.Errorf("Expected %d uploads, got %d", e, a)
	}
}
