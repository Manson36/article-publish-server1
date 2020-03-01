package datamodels

import "time"

type AdminUser struct {
	ID           int64        `json:"id, string"`
	NickName     string       `json:"nickName" gorm:"type:varchar(50);not null;default:'';pq_comment:用户昵称"`
	Email        string       `json:"email" gorm:"type:varchar(200);not null;unique;pq_comment:账号邮箱"`
	Password     string       `json:"password" gorm:"type:varchar(200);not null;pq_comment:用户密码"`
	Salt         string       `json:"salt" gorm:"type:char(5);not null;pq_comment:加密密码的盐"`
	AdminType    int          `json:"adminType" gorm:"type:smallint;not null;default:2;pq_comment:管理用户类型,1表示管理员，2表示普通管理者"`
	PlatformType PlatformType `json:"platformType" gorm:"type:smallint;pq_comment:管理员所处的平台类型"`
	CreatedAt    *time.Time   `json:"createdAt" gorm:"type:timestamptz;not null;default:now();pq_comment:该用户的创建时间"`
	UpdatedAt    *time.Time   `json:"updatedAt" gorm:"type:timestamptz;default:now();pq_comment:该用户信息的更新时间"`
	RemovedAt    *time.Time   `json:"removedAt" gorm:"type:timestamptz;pq_comment:用户被移除的时间"`
	Removed      bool         `json:"removed" gorm:"用户是否被移除"`
	Disabled     bool         `json:"disabled" gorm:"该用户是否被禁用"`
	DisabledAt   *time.Time   `json:"disabledAt" gorm:"type:timestamptz;pq_comment:用户的禁用时间"`
}
