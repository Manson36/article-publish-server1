package models

import "github.com/article-publish-server1/datamodels"

type ArticleCreateReqBody struct {
	Title     string                  `json:"title"`
	Author    string                  `json:"author"`
	Tags      []string                `json:"tags"`
	Summary   string                  `json:"summary"`
	Content   string                  `json:"content"`
	Cover     string                  `json:"cover"`
	CoverPath string                  `json:"coverPath"`
	Status    string                  `json:"status"` //1.表示草稿状态，2.表示发布状态
	Platform  datamodels.PlatformType `json:"platform"`
}

type ArticleInfoReqBody struct {
	ID       int64                   `json:"id, string"`
	Platform datamodels.PlatformType `json:"platform"`
}

type ArticleRemoveReqBody struct {
	ArticleInfoReqBody
}

type ArticleUpdateReqBody struct {
	ID int64 `json:"id, string"`
	ArticleCreateReqBody
}

type ArticleListReqBody struct {
	Page     PaginationCondition     `json:"page"`
	Status   int8                    `json:"status"`
	Platform datamodels.PlatformType `json:"platform"`
	Filter   string                  `json:"filter"`
}
