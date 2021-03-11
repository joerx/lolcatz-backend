package s3

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

var sess *session.Session

func init() {
	var err error
	sess, err = session.NewSession(&aws.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

// UploadRequest contains the parameters to process a file upload to S3
type UploadRequest struct {
	Filename     string
	OriginalName string
	ContentType  string
}

// Config struct with general settings for the upload
type Config struct {
	Region   string
	Bucket   string
	Endpoint string
}

func newClient(cf Config) *s3.S3 {
	awsCfg := aws.NewConfig().WithRegion(cf.Region).WithEndpoint(cf.Endpoint).WithS3ForcePathStyle(true)
	return s3.New(sess, awsCfg)
}

// Upload uploads a file to S3 base on given request
func Upload(r UploadRequest, cf Config) (string, error) {
	s3c := newClient(cf)
	up := s3manager.NewUploaderWithClient(s3c)

	file, err := os.Open(r.Filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	key := uuid.New().String()
	ext := filepath.Ext(r.OriginalName)
	if ext != "" {
		key += ext
	}

	tags := url.Values{}
	tags.Add("filename", r.OriginalName)

	input := &s3manager.UploadInput{
		Bucket:      aws.String(cf.Bucket),
		Key:         aws.String(key),
		Body:        file,
		Tagging:     aws.String(tags.Encode()),
		ContentType: aws.String(r.ContentType),
	}

	if _, err := up.Upload(input); err != nil {
		return "", fmt.Errorf("upload to S3 failed: %v", err)
	}

	s3url := fmt.Sprintf("s3://%s/%s", cf.Bucket, key)

	log.Printf("Successfully uploaded file to %s", s3url)

	return key, nil
}

// Presign generates a pre-signed URL for given S3 key
func Presign(key string, cf Config) (string, error) {
	svc := s3.New(sess, aws.NewConfig().WithRegion(cf.Region))

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(cf.Bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("Failed to generated pre-signed object request %v", err)
	}

	return urlStr, nil
}
