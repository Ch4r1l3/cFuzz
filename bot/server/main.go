package main

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/bot/server/router"
	"github.com/Ch4r1l3/cFuzz/bot/server/service"
)

func main() {

	config.Setup()
	models.Setup()
	defer models.DB.Close()
	service.Setup()
	defer service.StopFuzz()

	router := router.InitRouter()

	router.Run(fmt.Sprintf("0.0.0.0:%d", config.ServerConf.Port))
}
