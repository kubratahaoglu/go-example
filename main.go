package main

import (
	"mail-service/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", handler.PingGet)
	r.POST("/sendmail", handler.SendMailPost)
	r.Run()
}
