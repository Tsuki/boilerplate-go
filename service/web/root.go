package http

import (
	. "boilerplate-go/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/appleboy/gin-jwt"
	"time"
	"github.com/gin-contrib/sessions"
	"github.com/toorop/gin-logrus"
	"github.com/gin-contrib/sessions/memstore"
)

func StartGinServer() {
	app := gin.Default()
	store := memstore.NewStore([]byte("gDcNivjCoAbcTFTG6yra"))
	app.Use(sessions.Sessions("menSession", store))
	app.Use(ginlogrus.Logger(Log))

	app.GET("/", root)
	authMiddleware := jwtMiddleware()
	auth := app.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}
	app.Run(":8080")
}

func root(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"success": "OK"})
}

func jwtMiddleware() (*jwt.GinJWTMiddleware) {
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("xmloFDiEe5M0ttQv3irQ"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
				return userId, true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
	return authMiddleware
}
