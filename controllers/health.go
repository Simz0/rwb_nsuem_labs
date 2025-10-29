package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	message := "Server listening on port: " + os.Getenv("PORT")

	c.JSON(http.StatusOK, gin.H{"message": message})
}
