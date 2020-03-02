package services

import (
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/models"
	"github.com/article-publish-server1/utils/qnuploader"
	"mime/multipart"
)

type ImageService interface {
	UEImageUpload(file *multipart.FileHeader) map[string]interface{}
	ImageUptoken() *models.Ret
	ImageUEUptoken() *models.Ret
	ImageUploadCb(body *qnuploader.UploadImageCbBody) (*datamodels.Image, error)
	GetList(body *models.ImageListReqBody) *models.Ret
	RemoveImage(body *models.ImageRemoveReqBody) *models.Ret
}
