package main

import (
	"github.com/Dayasagara/aws-azure-replication/aws"
	"github.com/Dayasagara/aws-azure-replication/azure"
	"github.com/labstack/echo"
)
 
func main() {
	e := echo.New()
	e.POST("/replicate", replicate)
	e.Logger.Fatal(e.Start(":8095"))
}

func replicate(c echo.Context) error{
	aws.InsertIntoS3Bucket(c)
	azure.InsertIntoBLOBContainer(c)
	return nil
}

