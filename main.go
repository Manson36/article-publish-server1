package main

import (
	"github.com/article-publish-server1/config"
	"github.com/article-publish-server1/controllers"
	"github.com/article-publish-server1/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	if err := web(); err != nil {
		log.Fatal(err.Error())
	}
}

func web() error {
	switch config.Server.Mode {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "dev":
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, config.Server.Name+"pong")
	}).GET("/name", func(c *gin.Context) {
		c.String(http.StatusOK, config.Server.Name)
	})

	//admin_user
	adminUserController := controllers.AdminUserController{
		Service: services.NewAdminUserService(),
	}
	r.
		Group("/admin/user").
		POST("create", adminUserController.Create)

	return r.Run(":" + config.Web.Port)
}
