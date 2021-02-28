package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/joerx/lolcatz-backend/db"
	"github.com/joerx/lolcatz-backend/db/pg"
	"github.com/joerx/lolcatz-backend/http/handlers"
	"github.com/joerx/lolcatz-backend/http/middleware"
	"github.com/joerx/lolcatz-backend/http/routing"
	"github.com/joerx/lolcatz-backend/s3"

	_ "github.com/lib/pq"
)

type config struct {
	CorsAllowOrigin string
	S3              s3.Config
	DB              pg.Config
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
		log.Println("WARNING: Access-Control-Allow-Origin is set to '*' which should be used for development purposes only!")
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
	flag.StringVar(&cf.BindAddr, "bind", "localhost:3000", "Bind http server to this address")

	flag.Parse()

	return cf
}

func initDB(ctx context.Context, cf pg.Config) (db.DB, error) {
	pgdb, err := pg.NewClient(ctx, cf)
	if err != nil {
		return nil, err
	}

	if err := db.InitSchema(pgdb); err != nil {
		return nil, err
	}

	log.Println("Database connection initialized")
	return pgdb, nil
}

func main() {
	// parse config
	cf := parseFlags()
	checkCfg(cf)

	ctx := context.Background()

	db, err := initDB(ctx, cf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	us := pg.NewUploadService(db)

	uploadHandler := handlers.NewUpload(cf.S3, us)
	healthHandler := handlers.NewHealth(db)

	r := routing.NewRouter()
	r.Filter(middleware.Logging)
	r.Filter(middleware.CorsWithOrigin(cf.CorsAllowOrigin))

	// setup application routes
	r.Handle("/", handlers.Status)
	r.Handle("/upload", uploadHandler.CreateUpload)
	r.Handle("/list", uploadHandler.FindUploads)
	r.Handle("/health", healthHandler.Health)

	// start http server
	log.Printf("Starting server at %s", cf.BindAddr)
	http.ListenAndServe(cf.BindAddr, nil)
}
