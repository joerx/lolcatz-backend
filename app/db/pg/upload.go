package pg

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joerx/lolcatz-backend/db"
	"github.com/joerx/lolcatz-backend/upload"
)

// NewUploadService returns a new instance of UploadService
func NewUploadService(db db.DB) *UploadService {
	return &UploadService{db: db}
}

// UploadService implements the UploadService interface using postgres as storage backend
type UploadService struct {
	db db.DB
}

// CreateUpload creates an upload
func (s *UploadService) CreateUpload(ctx context.Context, u *upload.Upload) (*upload.Upload, error) {
	q := `INSERT INTO uploads(username, filename, s3key, timestamp) 
		  VALUES($1, $2, $3, $4)`

	u.Timestamp = time.Now()
	r, err := s.db.Exec(ctx, q, u.Username, u.Filename, u.S3Url, u.Timestamp)
	if err != nil {
		return nil, err
	}

	u.ID, err = r.LastInsertId()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// FindUploads allows to find uploads with optional filters
func (s *UploadService) FindUploads(ctx context.Context, f *upload.Filter) ([]*upload.Upload, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := f.Username; v != nil {
		where, args = append(where, "username = $1"), append(args, *v)
	}

	q := `SELECT id, username, filename, s3key, timestamp 
			FROM uploads WHERE ` + strings.Join(where, " AND ") + `
			ORDER BY id DESC`

	rows, err := s.db.Query(ctx, q, args...)
	if err != nil {
		log.Printf("pg query error - %v", err)
		return nil, fmt.Errorf("Database query error")
	}
	defer rows.Close()

	result := make([]*upload.Upload, 0)

	for rows.Next() {
		u := &upload.Upload{}
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Filename,
			&u.S3Url,
			&u.Timestamp,
		); err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	return result, nil
}
