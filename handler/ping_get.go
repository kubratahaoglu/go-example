package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//PingGet function
func PingGet(c *gin.Context) {

	c.JSON(http.StatusOK, map[string]string{
		"message": os.Getenv("HOSTNAME"),
	})

}
