package application

import (
	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/bookstore_utils-go/logger"
	"github.com/superbkibbles/realestate_users-api/src/datasources/mysqlclient"
)

var (
	router = gin.Default()
)

func StartApp() {
	mysqlclient.Init()
	mapUrls()
	logger.Info("App started")
	router.Static("assets", "datasources/images")
	router.Run(":8080")
}
