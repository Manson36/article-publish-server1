package models

type Ret struct {
	Code         int
	Msg          string
	TokenInvalid bool //如果是true，前端需要重新登录
	Data         interface{}
}
