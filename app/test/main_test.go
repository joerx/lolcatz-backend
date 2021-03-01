package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joerx/lolcatz-backend/db"
	"github.com/joerx/lolcatz-backend/db/pg"
	"github.com/joerx/lolcatz-backend/s3"
	"github.com/joerx/lolcatz-backend/server"
	"github.com/joerx/lolcatz-backend/test/testdb"

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

func initS3() (s3.Config, error) {
	// Create a new S3 bucket just for testing or use an existing one
	// Upload paths are randomized, so reusing an existin bucket should be fine
	// Creating a new bucket each time is safer, but what happens if the bucket
	// fails to delete? We'd see a lot of old buckets being created over time...
	return s3.Config{}, nil
}

func testMain(m *testing.M) int {
	var err error

	// Create database connection
	dbc, err = testdb.New(cfg)
	if err != nil {
		log.Printf("Error setting up database - %s", err)
		return 1
	}

	s3, err := initS3()
	if err != nil {
		log.Printf("Error setting up S3 - %s", err)
		return 1
	}

	cfg := server.DefaultConfig()
	cfg.S3 = s3

	// Create application instance
	app = server.New(cfg, dbc)

	return m.Run()
}

func Test_listUploads(t *testing.T) {
	// truncate database
	if err := testdb.Truncate(dbc); err != nil {
		t.Fatal(err)
	}

	// seed uploads
	if err := testdb.SeedUploads(dbc); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/list"), nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(body))

	if e, a := http.StatusOK, w.Result().StatusCode; e != a {
		t.Errorf("Expected statuscode %d but got %d", e, a)
	}
}
