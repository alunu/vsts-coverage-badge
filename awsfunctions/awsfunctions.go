package awsfunctions

import (
	"bytes"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// UploadToBucket uploads the given file to the filename and bucket
func UploadToBucket(bucket, filename, mimetype string, data *[]byte) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-2"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(*data),
		ContentType: aws.String(mimetype),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return fmt.Errorf("Unable to upload file: %v", err)
	}

	log.Println("Successfully uploaded file.")
	return nil
}
