package main

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/router"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"os"
)

func main() {
	config.Setup()
	os.MkdirAll(config.ServerConf.CrashesPath, os.ModePerm)
	logger.Setup()
	models.Setup()
	defer models.DB.Close()
	service.Setup()

	router := router.InitRouter()

	router.Run(fmt.Sprintf("0.0.0.0:%d", config.ServerConf.Port))

}
