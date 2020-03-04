package controllers

import (
	"github.com/article-publish-server1/models"
	"github.com/article-publish-server1/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TagController struct {
	Service services.TagService
}

func (t *TagController) Create(c *gin.Context) {
	body := models.TagCreateReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("添加标签失败，参数解析错误，errmsg：", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "添加标签失败，参数解析错误", Data: nil})
		return
	}

	c.JSON(http.StatusOK, t.Service.Create(&body))
}

func (t *TagController) Remove(c *gin.Context) {
	body := models.TagRemoveReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("删除标签失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "删除标签失败，参数解析错误", Data: nil})
		return
	}

	c.JSON(http.StatusOK, t.Service.Remove(&body))
}

func (t *TagController) Get(c *gin.Context) {
	body := models.TagInfoReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("获取标签信息失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "获取标签信息失败，参数解析错误", Data: nil})
		return
	}
	c.JSON(http.StatusOK, t.Service.Get(&body))
}

func (t *TagController) Update(c *gin.Context) {
	body := models.TagUpdateReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("更新标签信息失败，参数解析失败，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "更新标签信息失败，参数解析失败", Data: nil})
		return
	}

	c.JSON(http.StatusOK, t.Service.Update(&body))
}

func (t *TagController) ListAll(c *gin.Context) {
	body := models.TagListReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("标签列表获取失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "标签列表获取失败，参数解析错误", Data: nil})
		return
	}
	c.JSON(http.StatusOK, t.Service.ListAll(&body))
}
