package upload

import (
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Uploader interface {
	Exists(key string) (bool, error)
	Upload(key string, content io.ReadSeeker) error // Change content type to io.Reader
}

type S3UploaderImpl struct {
	client *s3.S3
	bucket string
}

var _ S3Uploader = &S3UploaderImpl{}

func NewS3Uploader(accessKey, secretKey, region, endpoint, bucket string) *S3UploaderImpl {
	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		log.Fatalf("Error creating new session: %v", err)
	}
	return &S3UploaderImpl{
		client: s3.New(sess),
		bucket: bucket,
	}
}

func (s *S3UploaderImpl) Exists(key string) (bool, error) {
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *S3UploaderImpl) Upload(key string, content io.ReadSeeker) error {
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   content, // Use the content as is, without converting to string
	})
	return err
}
