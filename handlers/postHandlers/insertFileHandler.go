package postHandlers

import (
	"log"
	"os"

	"github.com/Dayasagara/aws-azure-replication/aws"
	"github.com/Dayasagara/aws-azure-replication/azure"
	rm "github.com/Dayasagara/aws-azure-replication/responseManager"
	"github.com/labstack/echo"
)

type PostHandler struct{}

func (p *PostHandler) InsertHandler(ctx echo.Context) error {
	fileName := "file1"
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return rm.ResponseMapper(400, "File Not Found", ctx)
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	s3Err := make(chan error)
	blobErr := make(chan error)

	go aws.InsertIntoS3Bucket(file, s3Err, fileName)
	go azure.InsertIntoBLOBContainer(file, blobErr, fileName)

	s3Response := <-s3Err
	blobResponse := <-blobErr

	log.Println("s3: ", s3Response)
	log.Println("Az: ", blobResponse)

	if s3Response == nil && blobResponse == nil {
		return rm.ResponseMapper(200, "Successfully inserted", ctx)
	}

	return rm.ResponseMapper(400, "Error in inserting file", ctx)

}
