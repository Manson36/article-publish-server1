package services

import (
	"github.com/article-publish-server1/models"
	"github.com/article-publish-server1/utils/qnuploader"
)

type ArticleService interface {
	UploadCb(body *qnuploader.UploadImageCbBody) *models.Ret
	Uptoken() *models.Ret
	Create(body *models.ArticleCreateReqBody) *models.Ret
	Remove(body *models.ArticleRemoveReqBody) *models.Ret
	Get(body *models.ArticleInfoReqBody) *models.Ret
	Update(body *models.ArticleUpdateReqBody) *models.Ret
	List(body *models.ArticleListReqBody) *models.Ret
}
