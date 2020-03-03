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
		POST("/create", adminUserController.Create).
		POST("/login", adminUserController.Login).
		POST("/info", adminUserController.Info)

	//image
	imageControllser := controllers.ImageController{
		Service: services.NewImageService(),
	}
	r.
		Group("/image").
		GET("/ue", imageControllser.GetUEConfig).
		POST("/ue", imageControllser.UploadUEFile).
		POST("/uptoken", imageControllser.Uptoken).
		POST("/upload/db", imageControllser.UploadCb).
		POST("list", imageControllser.GetList).
		POST("/remove", imageControllser.Remove)

	//article
	articleController := controllers.ArticleController{
		Service: services.NewArticleService(),
	}
	r.
		Group("/article").
		POST("/uptoken", articleController.Uptoken).
		POST("/upload/cb", articleController.UploadCb).
		POST("/create", articleController.Create).
		POST("/remove", articleController.Remove).
		POST("/update", articleController.Update).
		POST("/info", articleController.Info).
		POST("/list", articleController.List)

	return r.Run(":" + config.Web.Port)
}
