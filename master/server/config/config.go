package config

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

type Server struct {
	RunMode      string `mapstructure:"runMode"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  string `mapstructure:"readTimeout"`
	WriteTimeout string `mapstructure:"writeTimeout"`
	TempPath     string `mapstructure:"tempPath"`
}

type Kubernetes struct {
	ConfigPath        string `mapstructure:"configPath"`
	Namespace         string `mapstructure:"namespace"`
	CheckTaskTime     int    `mapstructure:"checkTaskTime"`
	MaxClientRetryNum int    `mapstructure:"maxClientRetryNum"`
}

var ServerConf = &Server{}
var KubernetesConf = &Kubernetes{}

func Setup() {
	viper.SetConfigType("YAML")
	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal("Read 'config.yaml' fail: %v\n", err)
	}

	viper.ReadConfig(bytes.NewBuffer(data))
	viper.UnmarshalKey("server", ServerConf)
	viper.UnmarshalKey("kubernetes", KubernetesConf)
}