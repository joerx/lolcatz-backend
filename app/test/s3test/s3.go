package s3test

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joerx/lolcatz-backend/s3"
	"github.com/joerx/lolcatz-backend/util"

	s3api "github.com/aws/aws-sdk-go/service/s3"
)

const region = "ap-southeast-1"

var s3c *s3api.S3

func init() {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		log.Fatal(err)
	}
	s3c = s3api.New(sess, aws.NewConfig().WithRegion(region))
}

// Setup creates a new S3 bucket for the integration test
// Using a real bucket, we can be sure that the system behaves exactly like the real thing
// We can also transparently use something like localstack.cloud to make test cheaper and faster
func Setup() (*s3.Config, error) {
	region := "ap-southeast-1"
	endpoint := "http://localhost:4566" // localstack - FIXME: get this from env

	bucketName := fmt.Sprintf("lolcatzd-testbucket-%s", util.RandString(10))
	log.Printf("Test bucket %s", bucketName)

	cfg := &s3.Config{
		Bucket:   bucketName,
		Region:   region,
		Endpoint: endpoint,
	}

	if err := s3.MakeBucket(bucketName, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Teardown empties and deletes the test bucket
func Teardown(c *s3.Config) error {
	iter := s3manager.NewDeleteListIterator(s3c, &s3api.ListObjectsInput{
		Bucket: aws.String(c.Bucket),
	})
	if err := s3manager.NewBatchDeleteWithClient(s3c).Delete(aws.BackgroundContext(), iter); err != nil {
		return err
	}

	if _, err := s3c.DeleteBucket(&s3api.DeleteBucketInput{Bucket: &c.Bucket}); err != nil {
		return err
	}
	log.Printf("Deleted bucket %s", c.Bucket)

	return nil
}
