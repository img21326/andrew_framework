package helper

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Helper struct {
	client *s3.S3
	bucket string
	region string
}

func NewS3Helper(accessKey, secretKey, region, bucket string) *S3Helper {
	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}))
	return &S3Helper{
		client: s3.New(session),
		region: region,
		bucket: bucket,
	}
}

func (s *S3Helper) CreateFolder(folderName string, isPublic bool) error {
	config := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(folderName + "/"),
	}
	if isPublic {
		config.ACL = aws.String("public-read")
	}
	_, err := s.client.PutObject(config)
	return err
}

func (s *S3Helper) UploadFile(folderName, fileName string, fileBytes []byte) (string, error) {
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(folderName + "/" + fileName),
		Body:   bytes.NewReader(fileBytes),
	})
	url := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s/%s", s.bucket, s.region, folderName, fileName)
	return url, err
}
