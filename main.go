package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/joerx/lolcatz-backend/pkg/handlers"
	"github.com/joerx/lolcatz-backend/pkg/middleware"
	"github.com/joerx/lolcatz-backend/pkg/routing"
)

// const bucket = "sandbox-lolcatz-be-storage-468871832330"
// const region = "ap-southeast-1"

type config struct {
	CorsAllowOrigin string
	S3Region        string
	S3Bucket        string
}

func checkCfg(cf config) {
	if cf.S3Bucket == "" {
		log.Fatalf("bucket name is required")
	}
	if cf.S3Region == "" {
		log.Fatalf("AWS region is required")
	}
	if cf.CorsAllowOrigin == "*" {
		log.Println("WARNING: Access-Control-Allow-Origin is set to '*' which should be used for development purposes only!!!")
	}
}

func parseFlags() config {
	cf := config{}
	flag.StringVar(&cf.CorsAllowOrigin, "cors-allow-origin", "*", "Cors allow-origin header value")
	flag.StringVar(&cf.S3Bucket, "bucket", "", "S3 bucket to upload files to")
	flag.StringVar(&cf.S3Region, "region", "", "AWS region to connect to")

	flag.Parse()

	return cf
}

func main() {
	cf := parseFlags()
	checkCfg(cf)

	log.Println("Starting server")

	r := routing.NewRouter()
	r.Filter(middleware.Logging)
	r.Filter(middleware.CorsWithOrigin(cf.CorsAllowOrigin))

	r.Handle("/", handlers.Status)
	r.Handle("/upload", handlers.Upload(cf.S3Region, cf.S3Bucket))

	http.ListenAndServe("localhost:8000", nil)
}
