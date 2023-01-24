package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hi Back to you",
	})
}

func main() {
	application := gin.Default()
	application.GET("/ping", PingHandler)
	application.GET("/hello", HelloHandler)
	application.Run()
}
