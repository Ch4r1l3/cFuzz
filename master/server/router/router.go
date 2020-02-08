package router

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/controller"
	"github.com/Ch4r1l3/cFuzz/master/server/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(config.ServerConf.RunMode)
	userController := new(controller.UserController)
	api := r.Group("/api")
	api.Use(middleware.Auth)
	{
		deploymentController := new(controller.DeploymentController)
		api.GET("/deployment", middleware.Pagination, deploymentController.List)
		api.POST("/deployment", middleware.CheckUserExist, deploymentController.Create)
		api.PUT("/deployment/:id", middleware.CheckUserExist, deploymentController.Update)
		api.DELETE("/deployment/:id", middleware.CheckUserExist, deploymentController.Destroy)
		api.GET("/deployment/:path1", middleware.Pagination, controller.DeploymentGetHandler)

		taskController := new(controller.TaskController)
		api.GET("/task", middleware.Pagination, taskController.List)
		api.POST("/task", middleware.CheckUserExist, taskController.Create)
		api.POST("/task/:id/start", middleware.CheckUserExist, taskController.Start)
		api.POST("/task/:id/stop", middleware.CheckUserExist, taskController.Stop)
		api.PUT("/task/:id", middleware.CheckUserExist, taskController.Update)
		api.DELETE("/task/:id", middleware.CheckUserExist, taskController.Destroy)

		taskCrashController := new(controller.TaskCrashController)
		api.GET("/crash/:id", taskCrashController.Download)

		storageItemController := new(controller.StorageItemController)
		api.GET("/storage_item", middleware.Pagination, storageItemController.List)
		api.POST("/storage_item", middleware.CheckUserExist, storageItemController.Create)
		api.POST("/storage_item/exist", middleware.CheckUserExist, storageItemController.CreateExist)
		api.DELETE("/storage_item/:id", middleware.CheckUserExist, storageItemController.Destroy)
		api.GET("/storage_item/:path1", middleware.Pagination, controller.StorageItemGetHandler)

		api.GET("/task/:path1", controller.TaskGetHandler)
		api.GET("/task/:path1/:path2", middleware.Pagination, controller.TaskGetHandler)

		api.GET("/user/status", userController.Status)
		api.GET("/user", middleware.AdminOnly, middleware.Pagination, userController.List)
		api.POST("/user", middleware.AdminOnly, userController.Create)
		api.PUT("/user/:id", userController.Update)
		api.DELETE("/user/:id", middleware.AdminOnly, userController.Delete)

		docs := api.Group("docs")
		docs.Use(cors.Default())
		{
			docs.StaticFile("/swagger.json", "./docs/swagger.json")
		}
	}
	r.POST("/api/user/login", userController.Login)

	return r
}
