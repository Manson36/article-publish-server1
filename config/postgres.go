package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type postgres struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
	SSLMode  string
}

func (p *postgres) readConfig() {
	p.Host = viper.GetString("postgres.host")
	p.Port = viper.GetString("postgres.port")
	p.DB = viper.GetString("postgres.db")
	p.User = viper.GetString("postgres.user")
	p.Password = viper.GetString("postgres.password")
	p.SSLMode = viper.GetString("postgres.sslmode")
}

var Postgres = &postgres{}

func (p *postgres) GetURI() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DB, p.SSLMode)
}
