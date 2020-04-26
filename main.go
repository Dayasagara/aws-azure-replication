package main

import (
	"github.com/Dayasagara/aws-azure-replication/receivers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.POST("/replicate", receivers.Post.InsertHandler)
	e.Logger.Fatal(e.Start(":8095"))
}
