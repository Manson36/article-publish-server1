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
	imageController := controllers.ImageController{
		Service: services.NewImageService(),
	}
	r.
		Group("/image").
		GET("/ue", imageController.GetUEConfig).
		POST("/ue", imageController.UploadUEFile).
		POST("/uptoken", imageController.Uptoken).
		POST("/upload/db", imageController.UploadCb).
		POST("list", imageController.GetList).
		POST("/remove", imageController.Remove)

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

	//tag
	tagController := controllers.TagController{
		Service: services.NewTagService(),
	}

	r.
		Group("/tag").
		POST("/create", tagController.Create).
		POST("/remove", tagController.Remove).
		POST("/update", tagController.Update).
		POST("/info", tagController.Get).
		POST("/list", tagController.ListAll)

	return r.Run(":" + config.Web.Port)
}
