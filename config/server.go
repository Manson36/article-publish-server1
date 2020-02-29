package config

import "github.com/spf13/viper"

type server struct {
	Name string
	Mode string
}

func (s *server)readConf() {
	s.Name = viper.GetString("name")
	s.Mode = viper.GetString("mode")
}

var Server = &server{}
