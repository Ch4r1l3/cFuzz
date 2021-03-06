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
		imageController := new(controller.ImageController)
		api.GET("/image", middleware.Pagination, imageController.List)
		api.GET("/image/:id", imageController.Retrieve)
		api.POST("/image", middleware.CheckUserExist, imageController.Create)
		api.PUT("/image/:id", middleware.CheckUserExist, imageController.Update)
		api.DELETE("/image/:id", middleware.CheckUserExist, imageController.Destroy)

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
		api.GET("/storage_item/:type", middleware.Pagination, storageItemController.ListByType)
		api.POST("/storage_item", middleware.CheckUserExist, storageItemController.Create)
		api.POST("/storage_item/exist", middleware.CheckUserExist, storageItemController.CreateExist)
		api.DELETE("/storage_item/:id", middleware.CheckUserExist, storageItemController.Destroy)

		api.GET("/task/:path1", controller.TaskGetHandler)
		api.GET("/task/:path1/:path2", middleware.Pagination, controller.TaskGetHandler)

		api.GET("/user/info", userController.Info)
		api.GET("/user", middleware.AdminOnly, middleware.Pagination, userController.List)
		api.POST("/user", middleware.AdminOnly, userController.Create)
		api.PUT("/user/:id", userController.Update)
		api.DELETE("/user/:id", middleware.AdminOnly, userController.Delete)

	}
	r.POST("/api/user/login", userController.Login)
	r.GET("/api/user/logout", userController.Logout)
	docs := r.Group("api/docs")
	docs.Use(cors.Default())
	{
		docs.StaticFile("/swagger.json", "./docs/swagger.json")
	}

	return r
}
