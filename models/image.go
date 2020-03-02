package models

import "github.com/article-publish-server1/datamodels"

type ImageListReqBody struct {
	Page     PaginationCondition     `json:"page"`
	Platform datamodels.PlatformType `json:"platform"`
}

type ImageRemoveReqBody struct {
	ID       int64                   `json:"id, string"`
	Platform datamodels.PlatformType `json:"platform"`
}
