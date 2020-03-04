package models

import "github.com/article-publish-server1/datamodels"

type TagInfoReqBody struct {
	ID       int64                   `json:"id, string"`
	Platform datamodels.PlatformType `json:"platform"`
}

type TagRemoveReqBody struct {
	TagInfoReqBody
}

type TagCreateReqBody struct {
	Name     string                  `json:"name"`
	Platform datamodels.PlatformType `json:"platform"`
}

type TagUpdateReqBody struct {
	ID int64 `json:"id, string"`
	TagCreateReqBody
}

type TagListReqBody struct {
	Platform datamodels.PlatformType `json:"platform"`
}
