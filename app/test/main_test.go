package tests

import (
	"log"
	"net/http"
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

var pgCfg pg.Config = pg.Config{
	Host:     "localhost",
	Port:     5432,
	User:     "dbtest",
	Password: "dbtest",
	Name:     "dbtest",
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	var err error

	// Create database connection
	dbc, err = dbtest.New(pgCfg)
	if err != nil {
		log.Printf("Error setting up database - %s", err)
		return 1
	}
	defer dbc.Close()

	// Setup S3 bucket for testing
	s3cfg, err := s3test.Setup()
	if err != nil {
		log.Printf("Error setting up S3 - %s", err)
		return 1
	}
	defer s3test.Teardown(s3cfg)

	cfg := server.DefaultConfig()
	cfg.CorsAllowOrigin = corsOrigin
	cfg.S3 = s3cfg

	// Create application instance
	app = server.New(cfg, dbc)

	return m.Run()
}
