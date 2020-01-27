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
	r.POST("/task", taskController.Create)
	r.POST("/task/:id/start", taskController.Start)
	r.POST("/task/:id/stop", taskController.Stop)
	r.PUT("/task/:id", taskController.Update)
	r.DELETE("/task/:id", taskController.Destroy)

	taskCrashController := new(controller.TaskCrashController)
	r.GET("/crash/:id", taskCrashController.Download)

	storageItemController := new(controller.StorageItemController)
	r.GET("/storage_item", storageItemController.List)
	r.GET("/storage_item/:type", storageItemController.ListByType)
	r.POST("/storage_item", storageItemController.Create)
	r.POST("/storage_item/exist", storageItemController.CreateExist)
	r.DELETE("/storage_item/:id", storageItemController.Destroy)

	r.GET("/task/:path1", controller.TaskGetHandler)
	r.GET("/task/:path1/:path2", controller.TaskGetHandler)

	return r
}
