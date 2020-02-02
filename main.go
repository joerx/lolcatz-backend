package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/joerx/lolcatz-backend/pkg/db"
	"github.com/joerx/lolcatz-backend/pkg/handlers"
	"github.com/joerx/lolcatz-backend/pkg/middleware"
	"github.com/joerx/lolcatz-backend/pkg/routing"
	"github.com/joerx/lolcatz-backend/pkg/s3"

	_ "github.com/lib/pq"
)

// const bucket = "sandbox-lolcatz-be-storage-468871832330"
// const region = "ap-southeast-1"

// const(

// )

type config struct {
	CorsAllowOrigin string
	S3              s3.Config
	DB              db.Config
	BindAddr        string
}

func checkCfg(cf config) {
	if cf.S3.Bucket == "" {
		log.Fatalf("bucket name is required")
	}
	if cf.S3.Region == "" {
		log.Fatalf("AWS region is required")
	}
	if cf.CorsAllowOrigin == "*" {
		log.Println("WARNING: Access-Control-Allow-Origin is set to '*' which should be used for development purposes only!!!")
	}
}

func parseFlags() config {
	cf := config{}

	// cors flags
	flag.StringVar(&cf.CorsAllowOrigin, "cors-allow-origin", "*", "Cors allow-origin header value")

	// S3 flags
	flag.StringVar(&cf.S3.Bucket, "bucket", "", "S3 bucket to upload files to")
	flag.StringVar(&cf.S3.Region, "region", "", "AWS region to connect to")

	// database flags
	flag.StringVar(&cf.DB.Host, "db-host", "localhost", "Database connection host")
	flag.IntVar(&cf.DB.Port, "db-port", 5432, "Database connection host")
	flag.StringVar(&cf.DB.User, "db-user", "lolcatz", "Database user")
	flag.StringVar(&cf.DB.Password, "db-password", "default", "Database password")
	flag.StringVar(&cf.DB.Name, "db-name", "lolcatz", "Database name")

	// server vars
	flag.StringVar(&cf.BindAddr, "bind", "localhost:8000", "Bind http server to this address")

	flag.Parse()

	return cf
}

func initDB(cf db.Config) (*db.Client, error) {
	dbClient, err := db.NewClient(cf)
	if err != nil {
		return nil, err
	}

	if err := dbClient.InitSchema(); err != nil {
		return nil, err
	}

	log.Println("Database connection initialized")

	return dbClient, nil
}

func main() {
	// parse config
	cf := parseFlags()
	checkCfg(cf)

	// init database connection
	dbClient, err := initDB(cf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	// setup routing
	r := routing.NewRouter()
	r.Filter(middleware.Logging)
	r.Filter(middleware.CorsWithOrigin(cf.CorsAllowOrigin))

	r.Handle("/", handlers.Status)
	r.Handle("/upload", handlers.Upload(cf.S3, dbClient))
	r.Handle("/list", handlers.ListUploads(dbClient))

	// start http server
	log.Printf("Starting server at %s", cf.BindAddr)
	http.ListenAndServe(cf.BindAddr, nil)
}
