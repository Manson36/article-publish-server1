package models

import "github.com/article-publish-server1/datamodels"

type AdminUserAddReqBody struct {
	NickName     string                  `json:"nickName"` //用户昵称
	Email        string                  `json:"email"`    //用户登录的邮箱
	Password     string                  `json:"password"` //用户密码
	IsAdmin      bool                    `json:"isAdmin"`  //是否是管理员
	PlatformType datamodels.PlatformType `json:"platform"` //平台类型
}

type AdminUserLoginReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminUserLoginResBody struct {
	User  *datamodels.AdminUser `json:"user"`
	Token string                `json:"token"`
}
