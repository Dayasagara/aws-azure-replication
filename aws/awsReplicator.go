package aws

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func InsertIntoS3Bucket(file *os.File, s3Err chan error, fileName string) {
	err := godotenv.Load()
	if err != nil {
		s3Err <- err
		return
	}
	//Enter the name of the bucket
	bucket := os.Getenv("BUCKET")
	region := os.Getenv("REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		s3Err <- err
		return
	}

	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		s3Err <- err
		return
	}
	log.Printf("Successfully uploaded %q to %q\n", "file", bucket)

	s3Err <- err
}
