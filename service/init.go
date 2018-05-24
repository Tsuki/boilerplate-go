package service

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/toorop/gin-logrus"
	. "boilerplate-go/config"
)


func init() {
	Log.Info("Init Service")
	go startGinServer()
}

func startGinServer() {
	ginEngine := gin.Default()
	store := memstore.NewStore([]byte("gDcNivjCoAbcTFTG6yra"))
	ginEngine.Use(sessions.Sessions("menSession", store))
	ginEngine.Use(ginlogrus.Logger(Log))
	ginEngine.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"success": "OK"})
	})
	ginEngine.Run(":8080")
}
