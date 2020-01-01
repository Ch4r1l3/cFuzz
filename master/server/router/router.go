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

	taskController := new(controller.TaskController)
	r.GET("/task", taskController.List)
	r.POST("/task", taskController.Create)
	r.PUT("/task/:id", taskController.Update)
	r.DELETE("/task/:path1", controller.TaskDeleteHandler)

	taskCorpusController := new(controller.TaskCorpusController)
	r.GET("/task/:taskid/corpus", taskCorpusController.List)
	r.POST("/task/:taskid/corpus", taskCorpusController.Create)
	r.DELETE("/task/:path1/:path2/:path3", controller.TaskDeleteHandler)

	taskTargetController := new(controller.TaskTargetController)
	r.GET("/task/:taskid/target", taskTargetController.Retrieve)
	r.POST("/task/:taskid/target", taskTargetController.Create)
	r.DELETE("/task/:path1/:path2", controller.TaskDeleteHandler)

	fuzzerController := new(controller.FuzzerController)
	r.GET("/fuzzer", fuzzerController.List)
	r.POST("/fuzzer", fuzzerController.Create)
	r.DELETE("/fuzzer/:id", fuzzerController.Destroy)

	return r
}
