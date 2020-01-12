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

	deploymentController := new(controller.DeploymentController)
	r.GET("/deployment", deploymentController.List)
	r.POST("/deployment", deploymentController.Create)
	r.PUT("/deployment/:id", deploymentController.Update)
	r.DELETE("/deployment/:id", deploymentController.Destroy)

	taskController := new(controller.TaskController)
	r.GET("/task", taskController.List)
	r.GET("/task/:path1", controller.TaskGetHandler)
	r.POST("/task", taskController.Create)
	r.PUT("/task/:id", taskController.Update)
	r.DELETE("/task/:path1", controller.TaskDeleteHandler)

	taskCorpusController := new(controller.TaskCorpusController)
	r.GET("/task/:path1/:path2", controller.TaskGetHandler)
	r.POST("/task/:taskid/corpus", taskCorpusController.Create)
	r.DELETE("/task/:path1/:path2", controller.TaskDeleteHandler)

	taskTargetController := new(controller.TaskTargetController)
	r.POST("/task/:taskid/target", taskTargetController.Create)

	fuzzerController := new(controller.FuzzerController)
	r.GET("/fuzzer", fuzzerController.List)
	r.POST("/fuzzer", fuzzerController.Create)
	r.DELETE("/fuzzer/:id", fuzzerController.Destroy)

	return r
}
