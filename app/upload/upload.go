package upload

import (
	"context"
	"time"
)

// AssetSize is a logical representation of an artifacts dimensions
type AssetSize string

var Original AssetSize = "original"
var Thumbnail AssetSize = "thumbnail"
var Large AssetSize = "large"
var Medium AssetSize = "medium"
var Small AssetSize = "small"

// Upload represents an upload record stored in the database
type Upload struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Filename  string    `json:"filename"`
	S3Url     string    `json:"s3_url"`
	Timestamp time.Time `json:"timestamp"`
	Assets    []Asset
}

// Asset for an upload, i.e. a single image
// An upload can have multiple assets, e.g. representing the different sizes
type Asset struct {
	URL  string    `json:"url"`
	Size AssetSize `json:"size"`
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
