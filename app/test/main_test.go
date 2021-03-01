package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joerx/lolcatz-backend/db"
	"github.com/joerx/lolcatz-backend/db/pg"
	"github.com/joerx/lolcatz-backend/server"
	"github.com/joerx/lolcatz-backend/test/dbtest"
	"github.com/joerx/lolcatz-backend/test/s3test"

	_ "github.com/lib/pq"
)

var app http.Handler
var dbc db.DB

var cfg pg.Config = pg.Config{
	Host:     "localhost",
	Port:     5432,
	User:     "testdb",
	Password: "t3st",
	Name:     "testdb",
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	var err error

	// Create database connection
	dbc, err = dbtest.New(cfg)
	if err != nil {
		log.Printf("Error setting up database - %s", err)
		return 1
	}
	defer dbc.Close()

	// Setup S3 bucket for testing
	s3, err := s3test.Setup()
	if err != nil {
		log.Printf("Error setting up S3 - %s", err)
		return 1
	}
	defer s3test.Teardown(s3)

	cfg := server.DefaultConfig()
	cfg.S3 = s3

	// Create application instance
	app = server.New(cfg, dbc)

	return m.Run()
}

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
