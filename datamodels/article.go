package datamodels

import "time"

type Article struct {
	ID           int64        `json:"id, string" gorm:"type:bigint"`
	Title        string       `json:"title" gorm:"type:varchar(100);pq_comment:文章标题"`
	Summary      string       `json:"summary" gorm:"type:varchar(200);pq_comment:文章摘要"`
	Author       string       `json:"author" gorm:"type:varchar(100);pq_comment:作者"`
	Content      string       `json:"content" gorm:"type:text;pq_comment:文章内容，html代码"`
	Cover        int64        `json:"cover,string" gorm:"type:bigint;pq_comment:文章的封面图片，对应file表中的id"`
	CoverPath    string       `json:"coverPath" gorm:"type:varchar(200);pq_comment:文章封面图片的路径"`
	Uploader     int64        `json:"uplaoder" gorm:"type:bigint;pq_comment:上传者id，与admin_user表关联"`
	PlatformType PlatformType `json:"platform" gorm:"type:smallint;pq_comment:文章的平台类型"`
	Status       int8         `json:"status" gorm:"type:smallint;pq_comment:文章状态，1.表示草稿，2表示发布"`
	Tags         JsonNumArray `json:"tags" gorm:"type:jsonb;pq_comment:文章标签;default:'[]'::jsonb"`
	CreatedAt    *time.Time   `json:"createdAt" gorm:"type:timestamptz;default:now();pq_comment:文件的创建时间"`
	UpdatedAt    *time.Time   `json:"updatedAt" gorm:"type:timestamptz;default:now();pq_comment:文件的更新时间"`
	RemovedAt    *time.Time   `json:"removedAt" gorm:"type:timestamptz;default:null;pq_comment:文章的移除时间"`
	Remove       bool         `json:"remove" gorm:"type:boolean;default:FALSE;pq_comment:文章是否被移除"`
}
