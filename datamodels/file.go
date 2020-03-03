package datamodels

import "time"

type File struct {
	ID           int64        `json:"id, string"`
	Name         string       `json:"name" gorm:"type:varchar(200);pq_comment:文件名称"`
	FSize        int64        `json:"fsize" gorm:"type:bigint;pq_comment:文件大小"`
	Path         string       `json:"path" gorm:"type:varchar(200);pq_comment:文件路径"`
	Mime         string       `json:"mime" gorm:"type:varchar(100);pq_comment:文件的mime类型"`
	ExtName      string       `json:"extName" gorm:"type:varchar(10);pq_comment:文件的扩展名"`
	Hash         string       `json:"hash" gorm:"type:varchar(200);pq_comment:文件指纹"`
	PlatformType PlatformType `json:"platformType" gorm:"type:smallint;pq_comment:文件平台类型"`
	UploadedAt   *time.Time   `json:"uploadedAt" gorm:"type:timestamptz;default:now();pq_comment:文件上传时间"`
}
