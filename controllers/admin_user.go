package controllers

import (
	"github.com/article-publish-server1/models"
	"github.com/article-publish-server1/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminUserController struct {
	Service services.AdminUserService
}

func (a *AdminUserController) Create(c *gin.Context) {
	body := models.AdminUserAddReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "创建用户参数解析错误", Data: nil})
		return
	}

	c.JSON(http.StatusOK, a.Service.Create(&body))
}
