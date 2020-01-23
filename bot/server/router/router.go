package router

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(config.ServerConf.RunMode)
	r.MaxMultipartMemory = config.ServerConf.UploadFileLimit << 20

	storageItemController := new(controller.StorageItemController)
	r.GET("/storage_item", storageItemController.List)
	r.GET("/storage_item/:id", storageItemController.Retrieve)
	r.POST("/storage_item", storageItemController.Create)
	r.POST("/storage_item/exist", storageItemController.CreateExist)
	r.DELETE("/storage_item/:id", storageItemController.Destroy)

	taskController := new(controller.TaskController)
	r.GET("/task", taskController.Retrieve)
	r.POST("/task", taskController.Create)
	r.POST("/task/start", taskController.StartFuzz)
	r.POST("/task/stop", taskController.StopFuzz)
	r.DELETE("/task", taskController.Destroy)

	taskCrashController := new(controller.TaskCrashController)
	r.GET("/task/crash", taskCrashController.List)
	r.GET("/task/crash/:id", taskCrashController.Download)

	taskResultController := new(controller.TaskResultController)
	r.GET("/task/result", taskResultController.Retrieve)

	return r
}
