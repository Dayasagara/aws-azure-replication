package aws

import (
	"log"
	"os"
	"fmt"
	"github.com/labstack/echo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	rm "github.com/Dayasagara/aws-azure-replication/responseManager"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func InsertIntoS3Bucket(c echo.Context) error {
	//Enter the name of the bucket
	bucket := ""
	fileName := "test"
	region := ""
	file, err := os.Open(fileName)
	defer file.Close()
	
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key: aws.String(fileName),
		Body: file,
	})
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", "file", bucket, err)
	}
	log.Printf("Successfully uploaded %q to %q\n", "file", bucket)

	return rm.ResponseMapper(200, "Succesfully uploaded to AWS", c)

}
