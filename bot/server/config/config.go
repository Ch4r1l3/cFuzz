package config

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

type server struct {
	RunMode      string `mapstructure:"runMode"`
	Port         string `mapstructure:"port"`
	ReadTimeout  string `mapstructure:"readTimeout"`
	WriteTimeout string `mapstructure:"writeTimeout"`
}

var ServerConf = &server{}

func Setup() {
	viper.SetConfigType("YAML")
	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal("Read 'config.yaml' fail: %v\n", err)
	}

	viper.ReadConfig(bytes.NewBuffer(data))
	viper.UnmarshalKey("server", ServerConf)
}
