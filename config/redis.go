package config

import "github.com/spf13/viper"

type redis struct {
	Host      string
	Port      string
	DB        int
	Password  string
	KeyPrefix string
}

func (r *redis) readConfig() {
	r.Host = viper.GetString("redis.host")
	r.Port = viper.GetString("redis.port")
	r.DB = viper.GetInt("redis.db")
	r.Password = viper.GetString("redis.password")
	r.KeyPrefix = viper.GetString("redis.key_prefix")
}

var Redis = &redis{}
