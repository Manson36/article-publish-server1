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

type ImageController struct {
	Service services.ImageService
}

func (i *ImageController) GetUEConfig(c *gin.Context) {
	jsonStr := `{
  "imageActionName": "uploadimage",
  "imageFieldName": "upfile",
  "imageMaxSize": 4194304,
  "imageAllowFiles": [".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"],
  "imageCompressEnable": true,
  "imageCompressBorder": 1600,
  "imageInsertAlign": "none",
  "imageUrlPrefix": "",
  "imagePathFormat": "",
  "scrawlActionName": "uploadscraw",
  "scrawlFieldName": "upfile",
  "scrawlPathFormat": "",
  "scrawlMaxSize": 4194304,
  "scrawlUrlPrefix": "",
  "scrawlInsertAlign": "none"
}
`

	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, jsonStr)
}

func (i *ImageController) UploadUEFile(c *gin.Context) {
	uploadType := c.Query("action")
	if uploadType != "uploadimage" {
		c.JSON(http.StatusOK, gin.H{"state": "FAIL"})
		return
	}

	file, _ := c.FormFile("upfile")
	c.JSON(http.StatusOK, i.Service.UEImageUpload(file))
}

func (i *ImageController) Uptoken(c *gin.Context) {
	c.JSON(http.StatusOK, i.Service.ImageUptoken())
}

func (i *ImageController) UploadCb(c *gin.Context) {
	body := qnuploader.UploadImageCbBody{}
	buf, err := c.GetRawData()
	if err != nil {
		log.Println("图片素材上传回调失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "图片素材上传回调失败，参数解析错误", Data: err.Error()})
		return
	}

	buf = body.HandleNullString(buf)
	if err := json.Unmarshal(buf, &body); err != nil {
		log.Println("图片素材上传回调失败，参数解析错误，errmsg:", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "图片素材上传回调失败，参数解析错误", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, i.Service.ImageUploadCb)
}

func (i *ImageController) GetList(c *gin.Context) {
	body := models.ImageListReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("获取图片列表信息失败，参数解析错误, errmsg:", err.Error())
		c.JSON(http.StatusOK, &models.Ret{Code: 400, Msg: "获取图片列表信息失败，参数解析错误"})
		return
	}

	c.JSON(http.StatusOK, i.Service.GetList(&body))
}

func (i *ImageController) Remove(c *gin.Context) {
	body := models.ImageRemoveReqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("删除图片失败，参数解析错误，errmsg：", err.Error())
		c.JSON(http.StatusOK, models.Ret{Code: 400, Msg: "删除图片失败，参数解析错误"})
		return
	}

	c.JSON(http.StatusOK, i.Service.RemoveImage(&body))
}
