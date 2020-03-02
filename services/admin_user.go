package services

import (
	"github.com/article-publish-server1/datamodels"
	"github.com/article-publish-server1/models"
	"github.com/article-publish-server1/repositories"
	"github.com/article-publish-server1/utils"
	"log"
	"strings"
)

type AdminUserService interface {
	Create(body *models.AdminUserAddReqBody) *models.Ret
	Login(body *models.AdminUserLoginReqBody) *models.Ret
}

type adminUserService struct {
	repo repositories.AdminUserRepository
}

func NewAdminUserService() AdminUserService {
	repo := repositories.NewAdminUserRepository()
	return &adminUserService{
		repo: repo,
	}
}

func (a *adminUserService) Login(body *models.AdminUserLoginReqBody) *models.Ret {
	user, err := a.repo.Get("email = ? AND removed IS NOT TRUE", body.Email)
	if err != nil {
		log.Println("管理用戶信息获取失败", err.Error())
		return &models.Ret{Code: 500, Msg: "用户信息获取失败，请与平台联系"}
	}

	if user == nil {
		return &models.Ret{Code: 400, Msg: "该用户不存在"}
	}

	if utils.HashPwdWithSalt(body.Password, user.Salt) != user.Password {
		return &models.Ret{Code: 400, Msg: "输入密码错误"}
	}

	claims := models.AdminUserCustomClaims{}
	claims.UserID = user.ID
	token, err := claims.Sign()
	if err != nil {
		log.Println("管理用户登录token生成失败：", err.Error())
		return &models.Ret{Code: 500, Msg: "管理用户登录token生成失败"}
	}

	data := models.AdminUserLoginResBody{
		User:  user,
		Token: token,
	}

	return &models.Ret{Code: 200, Msg: "登录成功", Data: data}
}

func (a *adminUserService) Create(body *models.AdminUserAddReqBody) *models.Ret {
	id, err := utils.GetInt64ID()
	if err != nil {
		log.Println("创建账号时，获取生成的id错误:", err.Error())
		return &models.Ret{Code: 500, Msg: "创建账号时，生成id错误"}
	}

	switch body.PlatformType {
	case datamodels.ZingglobalPlatform, datamodels.ZhidreamPlatform, datamodels.HealthEnginePlatform:
	default:
		return &models.Ret{Code: 400, Msg: "请输入正确的平台类型"}
	}

	pwd := strings.TrimSpace(body.Password)
	if pwd == "" {
		return &models.Ret{Code: 400, Msg: "请输入管理员密码"}
	}

	email := strings.TrimSpace(body.Email)
	if email == "" {
		return &models.Ret{Code: 400, Msg: "请输入邮箱"}
	}

	nickName := strings.TrimSpace(body.NickName)
	if nickName == "" {
		return &models.Ret{Code: 400, Msg: "请输入昵称"}
	}

	pwdInfo := utils.GenPwdAndSalt(pwd)
	user := datamodels.AdminUser{
		ID:           id,
		NickName:     nickName,
		Email:        email,
		Password:     pwdInfo.Password,
		Salt:         pwdInfo.Salt,
		AdminType:    2,
		PlatformType: body.PlatformType,
	}

	if body.IsAdmin {
		user.AdminType = 1
	}

	u, err := a.repo.Get(`removed IS NOT TRUE AND email =? AND platform_type = ?`, email, body.PlatformType)
	if err != nil {
		log.Println("用户信息获取失败，数据库错误， errmsg:", err.Error())
		return &models.Ret{Code: 500, Msg: "y用户信息获取失败，请与平台联系"}
	}

	if u != nil {
		return &models.Ret{Code: 400, Msg: "该账号已经存在"}
	}

	if err := a.repo.Create(&user); err != nil {
		log.Println("创建用户错误：", err.Error())
		return &models.Ret{Code: 500, Msg: "创建用户错误"}
	}

	return &models.Ret{Code: 200, Msg: "用户创建成功", Data: user}
}
