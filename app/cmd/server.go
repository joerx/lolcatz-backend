package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/joerx/lolcatz-backend/db/pg"
	"github.com/joerx/lolcatz-backend/server"

	_ "github.com/lib/pq"
)

type flags struct {
	server  server.Config
	db      pg.Config
	address string
}

func checkCfg(cf flags) {
	if cf.server.S3.Bucket == "" {
		log.Fatalf("bucket name is required")
	}
	if cf.server.S3.Region == "" {
		log.Fatalf("AWS region is required")
	}
	if cf.server.CorsAllowOrigin == "*" {
		log.Println("WARNING: Access-Control-Allow-Origin is set to '*' which should be used for development purposes only!")
	}
}

func parseFlags() flags {
	cf := flags{}

	// cors flags
	flag.StringVar(&cf.server.CorsAllowOrigin, "cors-allow-origin", "*", "Cors allow-origin header value")

	// S3 flags
	flag.StringVar(&cf.server.S3.Bucket, "bucket", "", "S3 bucket to upload files to")
	flag.StringVar(&cf.server.S3.Region, "region", "", "AWS region to connect to")

	// database flags
	flag.StringVar(&cf.db.Host, "db-host", "localhost", "Database connection host")
	flag.IntVar(&cf.db.Port, "db-port", 5432, "Database connection host")
	flag.StringVar(&cf.db.User, "db-user", "lolcatz", "Database user")
	flag.StringVar(&cf.db.Name, "db-name", "lolcatz", "Database name")
	flag.StringVar(&cf.db.Password, "db-password", "s3cret", "Database password")

	// server vars
	flag.StringVar(&cf.address, "bind", "localhost:3000", "Bind http server to this address")
	flag.Parse()

	log.Printf("bucket %s", cf.server.S3.Bucket)
	log.Printf("region %s", cf.server.S3.Region)
	log.Printf("cors-allow-origin %s", cf.server.CorsAllowOrigin)
	log.Printf("db %s@%s:%d/%s", cf.db.User, cf.db.Host, cf.db.Port, cf.db.Name)

	checkCfg(cf)

	return cf
}

func main() {
	// parse config
	cfg := parseFlags()

	db, err := pg.NewWithSchema(context.Background(), cfg.db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := server.New(cfg.server, db)

	// start http server
	log.Printf("Starting server at %s", cfg.address)
	http.ListenAndServe(cfg.address, r)
}
