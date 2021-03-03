package dbtest

import (
	"context"
	"fmt"
	"time"

	"github.com/joerx/lolcatz-backend/db"
	"github.com/joerx/lolcatz-backend/db/pg"
	"github.com/joerx/lolcatz-backend/upload"
)

// New creates a new test database instance
func New(cfg pg.Config) (db.DB, error) {
	return pg.NewWithSchema(context.Background(), cfg)
}

// Truncate truncates the database for a new test
func Truncate(db db.DB) error {
	if _, err := db.Exec(context.Background(), "TRUNCATE TABLE uploads;"); err != nil {
		return fmt.Errorf("Error truncating database - %s", err)
	}
	return nil
}

// SeedUploads inserts a test fixture of uploads
func SeedUploads(db db.DB) ([]upload.Upload, error) {
	ctx := context.Background()
	uploads := make([]upload.Upload, 10)

	for i := range uploads {
		u := upload.Upload{
			Username:  "johndoe",
			Filename:  fmt.Sprintf("example_%d.png", i+1),
			S3Url:     fmt.Sprintf("s3://foo/bar/example_%d.png", i+1),
			Timestamp: time.Now(),
		}
		uploads[i] = u

		q := `INSERT INTO uploads(username, filename, s3key, timestamp) 
		VALUES($1, $2, $3, $4)`

		if _, err := db.Exec(ctx, q, u.Username, u.Filename, u.S3Url, u.Timestamp); err != nil {
			return nil, fmt.Errorf("Error inserting test upload %d - %s", i+1, err)
		}
	}

	return uploads, nil
}
