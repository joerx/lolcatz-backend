package s3

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

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
	Region string
	Bucket string
}

// Upload uploads a file to S3 base on given request
func Upload(r UploadRequest, cf Config) error {

	s3 := s3.New(sess, aws.NewConfig().WithRegion(cf.Region))
	up := s3manager.NewUploaderWithClient(s3)

	file, err := os.Open(r.Filename)
	if err != nil {
		return err
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
		return fmt.Errorf("upload to S3 failed: %v", err)
	}

	log.Printf("Successfully uploaded file to s3://%s/%s", cf.Bucket, key)

	return nil
}
