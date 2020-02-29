package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/go/bin/config")
	viper.AddConfigPath("../")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("config file read error, msg:%s", err))
	}

	Server.readConf()
	Web.readConfig()
	Postgres.readConfig()
	Redis.readConfig()
}
