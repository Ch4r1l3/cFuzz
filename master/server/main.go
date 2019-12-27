package main

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/router"
)

func main() {
	config.Setup()
	models.Setup()
	defer models.DB.Close()

	router := router.InitRouter()

	router.Run(fmt.Sprintf("0.0.0.0:%d", config.ServerConf.Port))

}
