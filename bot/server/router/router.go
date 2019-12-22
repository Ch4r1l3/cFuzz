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

	fuzzerController := new(controller.FuzzerController)
	r.GET("/fuzzer", fuzzerController.List)
	r.POST("/fuzzer", fuzzerController.Create)
	r.DELETE("/fuzzer/:name", fuzzerController.Destroy)

	taskController := new(controller.TaskController)
	r.GET("/task", taskController.Retrieve)
	r.POST("/task", taskController.Create)
	r.DELETE("/task", taskController.Destroy)

	taskCrashController := new(controller.TaskCrashController)
	r.GET("/task/crash", taskCrashController.List)

	taskCorpusController := new(controller.TaskCorpusController)
	r.POST("/task/corpus", taskCorpusController.Create)
	r.DELETE("/task/corpus", taskCorpusController.Destroy)

	return r
}
