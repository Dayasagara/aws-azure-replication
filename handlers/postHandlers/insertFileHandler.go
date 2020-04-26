package postHandlers

import (
	"log"
	"os"
	"time"

	"github.com/Dayasagara/aws-azure-replication/aws"
	"github.com/Dayasagara/aws-azure-replication/azure"
	rm "github.com/Dayasagara/aws-azure-replication/responseManager"
	"github.com/labstack/echo"
)

type PostHandler struct{}

func (p *PostHandler) InsertHandler(ctx echo.Context) error {

	//Get the form file from the request
	formFile, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	//Open the file which is of type multipart.File
	file, err := formFile.Open()
	if err != nil {
		return rm.ResponseMapper(400, "File Not Found", ctx)
	}
	defer file.Close()

	//Delete the temporary file created in the end
	defer func() {
		err = os.Remove(formFile.Filename)
		if err != nil {
			log.Println("File not deleted")
		}
	}()

	//Convert multipart.File to os.File
	dst, err := os.Create(formFile.Filename)
	if err != nil {
		return rm.ResponseMapper(400, "File Not Found", ctx)
	}
	defer dst.Close()

	//Assign file name
	fileName := formFile.Filename + " " + time.Now().Format("2006-01-02 15:04:05")

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	//Creating error channels
	s3Err := make(chan error)
	blobErr := make(chan error)

	//Go routines for concurrency
	go aws.InsertIntoS3Bucket(dst, s3Err, fileName)
	go azure.InsertIntoBLOBContainer(dst, blobErr, fileName)

	//Receive the errors through channels
	s3Response := <-s3Err
	blobResponse := <-blobErr

	if s3Response == nil && blobResponse == nil {
		return rm.ResponseMapper(200, "Successfully inserted", ctx)
	}

	log.Println("s3: ", s3Response)
	log.Println("Az: ", blobResponse)

	return rm.ResponseMapper(400, "Error in inserting file", ctx)

}
