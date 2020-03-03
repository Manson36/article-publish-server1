package controllers

import (
	"encoding/json"
	"github.com/article-publish-server1/models"
	"github.com/article-publish-server1/services"
	"github.com/article-publish-server1/utils/qnuploader"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ArticleController struct {
	Service services.ArticleService
}

func (a *ArticleController) Uptoken(c *gin.Context) {
	c.JSON(http.StatusOK, a.Service.Uptoken())
}

func (a *ArticleController) UploadCb(c *gin.Context) {
	body := qnuploader.UploadImageCbBody{}
	buf, err := c.GetRawData()

	if err != nil {
		log.Println("文章封面图片上传回调失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "文章封面图片上传回调失败，参数解析错误", Data: err.Error()})
	}

	//需要处理一下"null" -> null
	buf = body.HandleNullString(buf)
	if err := json.Unmarshal(buf, &body); err != nil {
		log.Println("文章封面图片上传回调失败，参数解析错误, errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "文章封面图片上传回调失败，参数解析错误", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, a.Service.UploadCb(&body))
}

func (a *ArticleController) Create(c *gin.Context) {
	body := models.ArticleCreateReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("创建文章失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "创建文章失败，参数解析错误", Data: nil})
		return
	}

	c.JSON(http.StatusOK, a.Service.Create(&body))
}

func (a *ArticleController) Info(c *gin.Context) {
	body := models.ArticleInfoReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("获取文章信息失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "获取文章信息失败，参数解析错误", Data: nil})
		return
	}

	c.JSON(http.StatusOK, a.Service.Get(&body))
}

func (a *ArticleController) Remove(c *gin.Context) {
	body := models.ArticleRemoveReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("删除文章信息失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "删除文章信息失败，参数解析错误", Data: nil})
		return
	}

	c.JSON(http.StatusOK, a.Service.Remove(&body))
}

func (a *ArticleController) Update(c *gin.Context) {
	body := models.ArticleUpdateReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("更新文章信息失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, &models.Ret{Code: 400, Msg: "更新文章信息失败，参数解析错误", Data: nil})
		return
	}

	c.JSON(http.StatusOK, a.Service.Update(&body))
}

func (a *ArticleController) List(c *gin.Context) {
	body := models.ArticleListReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("获取文章列表失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, &models.Ret{Code: 400, Msg: "获取文章列表失败，参数解析错误", Data: nil})
		return
	}

	c.JSON(http.StatusOK, a.Service.List(&body))
}
