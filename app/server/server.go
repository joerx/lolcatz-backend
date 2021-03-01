package server

import (
	"net/http"

	"github.com/joerx/lolcatz-backend/db"
	"github.com/joerx/lolcatz-backend/db/pg"
	"github.com/joerx/lolcatz-backend/http/handlers"
	"github.com/joerx/lolcatz-backend/http/middleware"
	"github.com/joerx/lolcatz-backend/http/routing"
	"github.com/joerx/lolcatz-backend/s3"
)

// Config is the overall application config
type Config struct {
	S3              *s3.Config
	CorsAllowOrigin string
}

// DefaultConfig returns the server configuration with some defaults
// Note: note all configuration flag may have meaningful defaults
func DefaultConfig() Config {
	return Config{
		CorsAllowOrigin: "*",
	}
}

// New creates a new server instance with the given database connection
func New(cfg Config, db db.DB) http.Handler {
	us := pg.NewUploadService(db)

	uploadHandler := handlers.NewUpload(cfg.S3, us)
	healthHandler := handlers.NewHealth(db)

	r := routing.NewRouter()
	r.Filter(middleware.Logging)
	r.Filter(middleware.CorsWithOrigin(cfg.CorsAllowOrigin))

	// setup application routes
	r.Handle("/", handlers.Status)
	r.Handle("/upload", uploadHandler.CreateUpload)
	r.Handle("/list", uploadHandler.FindUploads)
	r.Handle("/health", healthHandler.Health)

	return r
}
