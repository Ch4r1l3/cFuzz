package router

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(config.ServerConf.RunMode)

	dockerfileController := new(controller.DockerfileController)
	r.GET("/dockerfile", dockerfileController.List)
	r.POST("/dockerfile", dockerfileController.Create)
	r.PUT("/dockerfile/:id", dockerfileController.Update)
	r.DELETE("/dockerfile/:id", dockerfileController.Destroy)

	return r
}
