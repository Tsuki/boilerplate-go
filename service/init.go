package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var router *gin.Engine

func init() {
	go startGinServer()
}

func startGinServer() {
	//gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"success": "OK"})
	})
	router.Run(":8080")
}
