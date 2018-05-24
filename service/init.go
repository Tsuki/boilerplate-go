package service

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/toorop/gin-logrus"
	. "boilerplate-go/config"
)

var router = gin.Default()

func init() {
	Log.Info("Init Service")
	go startGinServer()
}

func startGinServer() {
	router.Use(ginlogrus.Logger(Logger))
	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"success": "OK"})
	})
	router.Run(":8080")
}
