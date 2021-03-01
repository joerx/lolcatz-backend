package upload

import (
	"context"
	"time"
)

// Upload represents an upload record stored in the database
type Upload struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Filename  string    `json:"filename"`
	S3Url     string    `json:"s3_url"`
	Timestamp time.Time `json:"timestamp"`
}

// Filter can be used in find operations to narrow down find results
type Filter struct {
	ID       *string `json:"id"`
	Username *string `json:"username"`
	Offset   *int    `json:"offset"`
	Limit    *int    `json:"limit"`
}

// Service is the common interface to find/create uploads. Uploads are
// immutable, so we can only create new ones and find existing ones
type Service interface {
	CreateUpload(ctx context.Context, u *Upload) (*Upload, error)
	FindUploads(ctx context.Context, f *Filter) ([]*Upload, error)
}
