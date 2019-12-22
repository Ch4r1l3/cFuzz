package main

import (
	//"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/bot/server/router"
)

func main() {

	config.Setup()
	models.Setup()

	router := router.InitRouter()

	//router.Run(fmt.Sprintf("0.0.0.0:%d", config.ServerConf.Port))
	router.Run("0.0.0.0:8888")
}
