package config

import (
	//"bytes"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/spf13/viper"
	//"io/ioutil"
	"log"
)

type Server struct {
	RunMode      string `mapstructure:"runMode"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  string `mapstructure:"readTimeout"`
	WriteTimeout string `mapstructure:"writeTimeout"`
	TempPath     string `mapstructure:"tempPath"`
	CrashesPath  string `mapstructure:"crashesPath"`
	LogToFile    bool   `mapstructure:"logToFile"`
	LogFileDir   string `mapstructure:"logFileDir"`
	SigningKey   string `mapstructure:"signingKey"`
}

type Kubernetes struct {
	ConfigPath        string `mapstructure:"configPath"`
	Namespace         string `mapstructure:"namespace"`
	CheckTaskTime     int    `mapstructure:"checkTaskTime"`
	MaxClientRetryNum int    `mapstructure:"maxClientRetryNum"`
	MaxStartTime      int64  `mapstructure:"maxStartTime"`
	RequestTimeout    int    `mapstructure:"requestTimeout"`
	InitCleanup       bool   `mapstructure:"initCleanup"`
}

var ServerConf = &Server{}
var KubernetesConf = &Kubernetes{}

func Setup() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Read 'config.yaml' fail: %v\n", err)
	}
	viper.UnmarshalKey("server", ServerConf)
	viper.UnmarshalKey("kubernetes", KubernetesConf)
	if ServerConf.SigningKey == "" {
		ServerConf.SigningKey, err = utils.RandomString(10)
		if err != nil {
			log.Fatal("random string error: %v\n", err)
		}
		viper.Set("server", ServerConf)
		err = viper.WriteConfig()
		if err != nil {
			log.Fatal("save signingKey in config error: %v\n", err)
		}
	}
}
