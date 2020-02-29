package config

import "github.com/spf13/viper"

type web struct {
	Host string
	Port string
	JWTSecret string
	TokenKey string
	TokenDomain string
	ExpiresAt int
}

func (w *web) readConfig() {
	w.Host = viper.GetString("web.host")
	w.Port = viper.GetString("web.port")
	w.JWTSecret = viper.GetString("web.jwt_secret")
	w.TokenKey = viper.GetString("web.token_key")
	w.TokenDomain = viper.GetString("token_domain")
	w.ExpiresAt = viper.GetInt("web.token_exp")
}

var Web = &web{}
