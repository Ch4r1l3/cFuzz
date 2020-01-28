package router

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/controller"
	"github.com/Ch4r1l3/cFuzz/master/server/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(config.ServerConf.RunMode)
	api := r.Group("/api")
	{
		deploymentController := new(controller.DeploymentController)
		api.GET("/deployment", middleware.Pagination, deploymentController.List)
		api.GET("/deployment/simplist", deploymentController.SimpList)
		api.POST("/deployment", deploymentController.Create)
		api.PUT("/deployment/:id", deploymentController.Update)
		api.DELETE("/deployment/:id", deploymentController.Destroy)

		taskController := new(controller.TaskController)
		api.GET("/task", middleware.Pagination, taskController.List)
		api.POST("/task", taskController.Create)
		api.POST("/task/:id/start", taskController.Start)
		api.POST("/task/:id/stop", taskController.Stop)
		api.PUT("/task/:id", taskController.Update)
		api.DELETE("/task/:id", taskController.Destroy)

		taskCrashController := new(controller.TaskCrashController)
		api.GET("/crash/:id", taskCrashController.Download)

		storageItemController := new(controller.StorageItemController)
		api.GET("/storage_item", middleware.Pagination, storageItemController.List)
		api.GET("/storage_item/:type", middleware.Pagination, storageItemController.ListByType)
		api.POST("/storage_item", storageItemController.Create)
		api.POST("/storage_item/exist", storageItemController.CreateExist)
		api.DELETE("/storage_item/:id", storageItemController.Destroy)

		api.GET("/task/:path1", controller.TaskGetHandler)
		api.GET("/task/:path1/:path2", middleware.Pagination, controller.TaskGetHandler)
	}

	return r
}
