package services

import (
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/datasorces"
	"github.com/article-publish-server1/models"
	"github.com/article-publish-server1/repositories"
	"github.com/article-publish-server1/utils"
	"github.com/article-publish-server1/utils/qnuploader"
	"log"
	"mime/multipart"
	"sync"
)

type ImageService interface {
	ArticleCoverUptoken() *models.Ret
	UEImageUpload(file *multipart.FileHeader) map[string]interface{}
	ImageUptoken() *models.Ret
	ImageUEUptoken() *models.Ret
	ImageUploadCb(body *qnuploader.UploadImageCbBody) *models.Ret
	CreateImageByUploadBody(body *qnuploader.UploadImageCbBody) (*datamodels.Image, *models.Ret)
	GetList(body *models.ImageListReqBody) *models.Ret
	RemoveImage(body *models.ImageRemoveReqBody) *models.Ret
}

type imageService struct {
	uploader *qnuploader.Uploader
	repo     repositories.ImageRepository
	fileRepo repositories.FileRepository
}

func NewImageService() ImageService {
	uploader := qnuploader.NewUploader(nil)
	return &imageService{
		uploader: uploader,
		repo:     repositories.NewImageRepository(),
		fileRepo: repositories.NewFileRepository(),
	}
}

func (i imageService) ArticleCoverUptoken() *models.Ret {
	return &models.Ret{
		Code: 200,
		Msg:  "获取文章封面图片上传拼争成功",
		Data: map[string]string{
			"uptoken": i.uploader.GetImageUptoken(
				i.uploader.CallbackURI+"article/upload/cb", nil),
		},
	}
}

func (i imageService) UEImageUpload(file *multipart.FileHeader) map[string]interface{} {
	uptoken := i.uploader.GetImageUptoken("", nil)
	body := qnuploader.UploadImageCbBody{}

	key := "ow/platform/image/" + utils.GetRandomString(32)
	params := map[string]string{
		"x:uploader": "1",
		"x:filenaem": file.Filename,
		"x:platform": "1",
	}

	if err := i.uploader.UploadFormFile(&body, params, uptoken, key, file); err != nil {
		log.Println("编辑器上传图片失败，errmsg:", err.Error())
		return map[string]interface{}{"state": "FAIL", "msg": "编辑器上传图片失败，请与平台联系"}
	}
	img, ret := i.CreateImageByUploadBody(&body)
	if ret != nil {
		return map[string]interface{}{"state": "FAIL", "msg": ret.Msg}
	}

	return map[string]interface{}{"state": "SUCCESS", "url": img.Path, "title": img.Name, "original": img.Name}
}

func (i imageService) ImageUptoken() *models.Ret {
	return &models.Ret{
		Code: 200,
		Msg:  "获取图片上传凭证成功",
		Data: map[string]string{
			"uptoken": i.uploader.GetImageUptoken(
				i.uploader.CallbackURI+"image/upload/cb", nil),
		},
	}
}

func (i imageService) ImageUEUptoken() *models.Ret {
	return &models.Ret{
		Code: 200,
		Msg:  "获取UE图片上传凭证成功",
		Data: map[string]string{
			"uptoken": i.uploader.GetImageUptoken(
				i.uploader.CallbackURI+"image/upload/ue/cb", nil),
		},
	}
}

func (i imageService) ImageUploadCb(body *qnuploader.UploadImageCbBody) *models.Ret {
	ret, e := i.CreateImageByUploadBody(body)
	if e != nil {
		return e
	}

	return &models.Ret{Code: 200, Msg: "图片素材上传成功", Data: *ret}
}

func (i imageService) CreateImageByUploadBody(body *qnuploader.UploadImageCbBody) (*datamodels.Image, *models.Ret) {
	var (
		errRet       models.Ret
		img          *datamodels.Image
		err          error
		imageID      int64
		path         string
		platformType datamodels.PlatformType
	)

	imageFile, err := i.fileRepo.Get("hash = ?", body.Hash)
	if err != nil {
		log.Println("上传图片时，读取文件Hash错误，请与平台联系")
		errRet = models.Ret{Code: 501, Msg: "上传图片时，读取文件Hash错误，请与平台联系"}
		goto ERRRET
	}

	platformType = datamodels.PlatformType(body.Platform)
	switch platformType {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		errRet = models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
		goto ERRRET
	}

	imageID, err = utils.GetInt64ID()
	if err != nil {
		log.Println("保存图片时，获取图片id信息错误：", err.Error())
		errRet = models.Ret{Code: 500, Msg: "保存图片信息时，获取图片id信息错误，请与平台联系"}
		goto ERRRET
	}

	//图片的类型，上传者信息
	img = &datamodels.Image{
		ID:           imageID,
		Name:         body.Name,
		Width:        body.Width,
		Height:       body.Height,
		Uploader:     body.Uploader,
		PlatformType: platformType,
	}

	//文件已存在
	if imageFile != nil {
		path = imageFile.Path
		if imageFile.Path != body.Key {
			if err := i.uploader.DeleteBucketRes("static", body.Key); err != nil {
				log.Println("删除重复文件失败，errmsg:", err.Error())
			}
		}
		img.Link = imageFile.ID
		img.Path = imageFile.Path
		if err := i.repo.Create(img); err != nil {
			log.Println("保存图片信息失败，请与平台联系， errmsg：", err.Error())
			errRet = models.Ret{Code: 501, Msg: "保存图片信息失败，请与平台联系"}
			goto ERRRET
		}
	} else {
		path = body.Key
		fileID, err := utils.GetInt64ID()
		if err != nil {
			log.Println("保存图片文件时，获取图片文件id信息错误：", err.Error())
			errRet = models.Ret{Code: 500, Msg: "保存图片文件时，获取图片文件id信息错误，请与平台联系"}
			goto ERRRET
		}

		tx := datasorces.PqDB.Begin()
		if tx.Error != nil {
			errRet = models.Ret{Code: 501, Msg: "保存图片时，开始数据库事务失败，请与平台联系"}
			goto ERRRET
		}

		file := datamodels.File{
			ID:           fileID,
			Name:         body.Name,
			FSize:        body.Fsize,
			Path:         body.Key,
			Mime:         body.MimeType,
			ExtName:      body.ExtName,
			Hash:         body.Hash,
			PlatformType: platformType,
		}

		img.Link = fileID
		img.Path = body.Key

		if err := i.fileRepo.CreateWithTx(tx, &file); err != nil {
			tx.Rollback()
			log.Println("保存图片文件失败，errmsg：", err.Error())
			errRet = models.Ret{Code: 501, Msg: "保存图片文件失败，请与平台联系"}
			goto ERRRET
		}

		if err := i.repo.CreateWithTx(tx, img); err != nil {
			tx.Rollback()
			log.Println("保存图片信息失败，errmsg:", err.Error())
			errRet = models.Ret{Code: 501, Msg: "保存图片信息失败, 请与平台联系"}
			goto ERRRET
		}

		if err := tx.Commit().Error; err != nil {
			log.Println("图片保存失败，errmsg：", err.Error())
			errRet = models.Ret{Code: 501, Msg: "保存图片失败，请与平台联系"}
			goto ERRRET
		}
	}

	img.Path = i.uploader.StaticURI + path
	return img, nil

ERRRET:
	if body.Key != "" {
		if err := i.uploader.DeleteBucketRes("static", body.Key); err != nil {
			log.Println("业务错误，删除已上传文件失败，errmsg:", err.Error())
		}
	}
	return nil, &errRet
}

func (i imageService) GetList(body *models.ImageListReqBody) *models.Ret {
	var (
		wg       sync.WaitGroup
		total    int64
		list     []datamodels.Image
		totalErr error
		listErr  error
	)

	switch body.Platform {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		return &models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		list, listErr = i.repo.List(
			`created_at DESC`, body.Page.Limit, body.Page.Offset,
			`removed IS NOT TRUE AND platform_type = ?`, body.Platform)
	}()

	go func() {
		defer wg.Done()
		total, totalErr = i.repo.Count(`removed IS NOT TRUE AND platform = ?`, &body)
	}()

	wg.Wait()

	if listErr != nil {
		log.Println("获取图片列表失败，数据库错误，errmsg:", listErr.Error())
		return &models.Ret{Code: 500, Msg: "获取图片列表失败，请与平台联系"}
	}

	if totalErr != nil {
		log.Println("获取图片列表数量失败，数据库错误，errmsg:", totalErr.Error())
		return &models.Ret{Code: 500, Msg: "获取图片列表失败，请与平台联系"}
	}

	for j := len(list) - 1; j >= 0; j-- {
		list[j].Path = i.uploader.StaticURI + list[j].Path
	}

	return &models.Ret{Code: 200, Msg: "获取图片列表成功", Data: map[string]interface{}{
		"total": total,
		"list":  list,
	}}
}

func (i imageService) RemoveImage(body *models.ImageRemoveReqBody) *models.Ret {
	if body.ID == 0 {
		return &models.Ret{Code: 400, Msg: "请传入正确的图片唯一标识信息"}
	}

	switch body.Platform {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		return &models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
	}

	if err := i.repo.Remove("id = ? AND platform_type = ?", body.ID, body.Platform); err != nil {
		log.Println("图片删除失败，数据库错误，errmsg:", err.Error())
		return &models.Ret{Code: 500, Msg: "图片删除失败，请与平台联系"}
	}

	return &models.Ret{Code: 200, Msg: "图片删除成功"}
}
