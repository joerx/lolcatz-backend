package testdb

import (
	"context"
	"fmt"
	"time"

	"github.com/joerx/lolcatz-backend/db"
	"github.com/joerx/lolcatz-backend/db/pg"
)

func New(cfg pg.Config) (db.DB, error) {
	return pg.NewClientWithSchema(context.Background(), cfg)
}

func Truncate(db db.DB) error {
	if _, err := db.Exec(context.Background(), "TRUNCATE TABLE uploads;"); err != nil {
		return fmt.Errorf("Error truncating database - %s", err)
	}
	return nil
}

func SeedUploads(db db.DB) error {
	ctx := context.Background()

	for i := 1; i <= 10; i++ {
		username := "johndoe"
		filename := fmt.Sprintf("example_%d.png", i)
		s3Url := fmt.Sprintf("s3://foo/bar/example_%d.png", i)
		timestamp := time.Now()

		q := `INSERT INTO uploads(username, filename, s3key, timestamp) 
		VALUES($1, $2, $3, $4)`

		if _, err := db.Exec(ctx, q, username, filename, s3Url, timestamp); err != nil {
			return fmt.Errorf("Error loading test record %d - %s", i, err)
		}
	}

	return nil
}
