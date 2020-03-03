package services

import (
	"fmt"
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/models"
	"github.com/article-publish-server1/repositories"
	"github.com/article-publish-server1/utils"
	"github.com/article-publish-server1/utils/qnuploader"
	"log"
	"strings"
	"sync"
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

type articleService struct {
	imageSvc ImageService
	repo     repositories.ArticleRepository
	fileRepo repositories.FileRepository
	uploader *qnuploader.Uploader
}

func NewArticleService() ArticleService {
	uploader := qnuploader.NewUploader(nil)
	return &articleService{
		imageSvc: NewImageService(),
		repo:     repositories.NewArticleRepository(),
		fileRepo: repositories.NewFileRepository(),
		uploader: uploader,
	}
}

func (a articleService) UploadCb(body *qnuploader.UploadImageCbBody) *models.Ret {
	ret, e := a.imageSvc.CreateImageByUploadBody(body)
	if e != nil {
		return e
	}

	return &models.Ret{
		Code: 200,
		Msg:  "文章封面图片上传成功",
		Data: *ret,
	}
}

func (a articleService) Uptoken() *models.Ret {
	return a.imageSvc.ArticleCoverUptoken()
}

func (a articleService) Create(body *models.ArticleCreateReqBody) *models.Ret {
	tags, err := utils.StringSliceToJsonNumArray(body.Tags)
	if err != nil {
		return &models.Ret{Code: 400, Msg: "标签信息错误"}
	}

	switch body.Platform {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		return &models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
	}

	article := datamodels.Article{
		Title:        strings.TrimSpace(body.Title),
		Author:       strings.TrimSpace(body.Author),
		Summary:      strings.TrimSpace(body.Summary),
		Tags:         tags,
		Content:      strings.TrimSpace(body.Content),
		Cover:        body.Cover,
		Status:       body.Status,
		PlatformType: body.Platform,
	}

	if article.Title == "" {
		return &models.Ret{Code: 400, Msg: "请填写文章标题"}
	}

	if article.Status != 1 && article.Status != 2 {
		return &models.Ret{Code: 400, Msg: "请设置正确的文章状态"}
	}

	if article.Cover != 0 {
		file, err := a.fileRepo.Get("id = ?", article.Cover)
		if err != nil {
			log.Println("文章封面信息获取失败，数据库错误，errmsg:", err.Error())
			return &models.Ret{Code: 500, Msg: "文章封面信息获取失败，请与平台联系"}
		}

		if file == nil {
			return &models.Ret{Code: 400, Msg: "封面不存在，请重新上传"}
		}

		article.CoverPath = file.Path
	}

	if article.Status == 2 {
		if article.Content == "" {
			return &models.Ret{Code: 400, Msg: "请填写文章内容"}
		}

		if article.Cover == 0 || article.CoverPath == "" {
			return &models.Ret{Code: 400, Msg: "请上传文章封面图片"}
		}
	}

	id, err := utils.GetInt64ID()
	if err != nil {
		log.Println("新建文章时，获取id信息错误，errmsg:", err.Error())
		return &models.Ret{Code: 501, Msg: "新建文章时，获取id信息错误，请与平台联系"}
	}
	article.ID = id

	if err := a.repo.Create(&article); err != nil {
		log.Println("文章创建失败，数据库错误，errmsg:", err.Error())
		return &models.Ret{Code: 500, Msg: "文章创建失败，请与平台联系"}
	}

	if article.CoverPath != "" {
		article.CoverPath = a.uploader.StaticURI + article.CoverPath
	}

	if article.Status == 1 {
		return &models.Ret{Code: 200, Msg: "文章保存成功", Data: article}
	}

	return &models.Ret{Code: 200, Msg: "文章发布成功", Data: article}
}

func (a articleService) Remove(body *models.ArticleRemoveReqBody) *models.Ret {
	if body.ID == 0 {
		return &models.Ret{Code: 400, Msg: "请传入正确的文章唯一标识信息"}
	}

	switch body.Platform {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		return &models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
	}

	if err := a.repo.Remove("id = ? AND platform_type = ?", body.ID, body.Platform); err != nil {
		log.Println("文章删除失败，数据库错误，errmsg:", err.Error())
		return &models.Ret{Code: 500, Msg: "文章删除失败，请与平台联系"}
	}

	return &models.Ret{Code: 200, Msg: "文章删除成功"}
}

func (a articleService) Get(body *models.ArticleInfoReqBody) *models.Ret {
	if body.ID == 0 {
		return &models.Ret{Code: 400, Msg: "请传入正确的文章唯一标识信息"}
	}

	switch body.Platform {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		return &models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
	}

	article, err := a.repo.Get("id=? AND removed IS NOT TRUE AND platform_type = ?", body.ID, body.Platform)
	if err != nil {
		log.Println("获取文章信息失败，数据库错误，errmsg：", err.Error())
		return &models.Ret{Code: 500, Msg: "获取文章信息失败，请与平台联系"}
	}

	if article == nil {
		return &models.Ret{Code: 400, Msg: "文章信息不存在"}
	}

	if article.CoverPath != "" {
		article.CoverPath = a.uploader.StaticURI + article.CoverPath
	}

	return &models.Ret{Code: 200, Msg: "文章信息获取成功", Data: *article}
}

func (a articleService) Update(body *models.ArticleUpdateReqBody) *models.Ret {
	switch body.Platform {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		return &models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
	}
	article, err := a.repo.Get("id=? AND removed IS NOT TRUE AND platform_type = ?", body.ID, body.Platform)
	if err != nil {
		log.Println("获取文章信息失败，数据库错误，errmsg:", err.Error())
		return &models.Ret{Code: 500, Msg: "获取文章信息失败，数据库错误"}
	}

	if article == nil {
		return &models.Ret{Code: 400, Msg: "改文章不存在"}
	}

	article.Title = strings.TrimSpace(body.Title)
	article.Author = strings.TrimSpace(body.Author)
	article.Summary = strings.TrimSpace(body.Summary)
	article.Content = strings.TrimSpace(body.Content)
	article.Cover = body.Cover
	article.Status = body.Status

	tags, err := utils.StringSliceToJsonNumArray(body.Tags)
	if err != nil {
		return &models.Ret{Code: 400, Msg: "标签信息错误"}
	}
	article.Tags = tags

	if article.Title == "" {
		return &models.Ret{Code: 400, Msg: "请填写文章标题"}
	}

	if article.Status != 1 && article.Status != 2 {
		return &models.Ret{Code: 400, Msg: "请设置正确的文章状态"}
	}

	if article.Cover != 0 {
		file, err := a.fileRepo.Get("id = ?", article.Cover)
		if err != nil {
			log.Println("文章封面信息获取失败，数据库错误，errmsg:", err.Error())
			return &models.Ret{Code: 500, Msg: "文章封面信息获取失败,请与平台联系"}
		}

		if file == nil {
			return &models.Ret{Code: 400, Msg: "封面不存在，请重新上传"}
		}

		article.CoverPath = file.Path
	}

	if article.Status == 2 {
		if article.Content == "" {
			return &models.Ret{Code: 400, Msg: "请填写文章内容"}
		}

		if article.Cover == 0 {
			return &models.Ret{Code: 400, Msg: "请上传文章封面图片"}
		}
	}

	if err := a.repo.Save(article); err != nil {
		log.Println("文章信息更新失败，数据库错误，errmsg:", err.Error())
		return &models.Ret{Code: 500, Msg: "文章信息更新失败，请与平台联系"}
	}

	if article.CoverPath != "" {
		article.CoverPath = a.uploader.StaticURI + article.CoverPath
	}

	if article.Status == 1 {
		return &models.Ret{Code: 200, Msg: "文章更新成功", Data: *article}
	}

	return &models.Ret{Code: 200, Msg: "文章发布成功", Data: *article}
}

func (a articleService) List(body *models.ArticleListReqBody) *models.Ret {
	var (
		wg       sync.WaitGroup
		list     []datamodels.Article
		total    int64
		listErr  error
		totalErr error
	)

	switch body.Platform {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		return &models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
	}

	wg.Add(2)

	query := fmt.Sprintf(`removed IS NOT TRUE AND platform_type =%d`, body.Platform)
	if body.Filter != "" {
		query += fmt.Sprintf(`AND (title LIKE '%%%s%%' OR summary LIKE '%%%[1]s%%' OR author LIKE '%%%[1]s%%' OR content LIKE '%%%[1]s%%')`,
			body.Filter)
	}

	if body.Status != 0 {
		query += fmt.Sprintf(" AND status = %d", body.Status)
	}

	go func() {
		defer wg.Done()
		total, totalErr = a.repo.Count(query)
	}()

	wg.Wait()

	if listErr != nil {
		log.Println("获取文章列表失败，数据库错误，errmsg:", listErr.Error())
		return &models.Ret{Code: 500, Msg: "获取新闻列表失败，请与平台联系"}
	}

	if totalErr != nil {
		log.Println("获取文章列表数量失败，数据库错误，errmsg：", totalErr.Error())
		return &models.Ret{Code: 500, Msg: "获取文章列表数量失败，请与平台联系"}
	}

	for i := len(list) - 1; i >= 0; i-- {
		list[i].CoverPath = a.uploader.StaticURI + list[i].CoverPath
	}

	return &models.Ret{Code: 200, Msg: "获取文章列表成功", Data: map[string]interface{}{
		"total": total,
		"list":  list,
	}}
}
