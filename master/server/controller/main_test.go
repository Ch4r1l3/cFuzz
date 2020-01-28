package controller

import (
	"bytes"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/middleware"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	r *gin.Engine
)

func prepareTestDB() {
	var err error
	models.DB, err = gorm.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}

	models.DB.SingularTable(true)
	models.DB.AutoMigrate(&models.Deployment{}, &models.Task{}, &models.StorageItem{}, &models.TaskEnvironment{}, &models.TaskArgument{}, &models.TaskCrash{}, &models.TaskFuzzResult{}, &models.TaskFuzzResultStat{})

}

func prepareRouter() {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode("debug")

	api := r.Group("/api")
	{
		deploymentController := new(DeploymentController)
		api.GET("/deployment", middleware.Pagination, deploymentController.List)
		api.POST("/deployment", deploymentController.Create)
		api.PUT("/deployment/:id", deploymentController.Update)
		api.DELETE("/deployment/:id", deploymentController.Destroy)
		api.GET("/deployment/:path1", DeploymentGetHandler)

		taskController := new(TaskController)
		api.GET("/task", middleware.Pagination, taskController.List)
		api.POST("/task", taskController.Create)
		api.POST("/task/:id/start", taskController.Start)
		api.POST("/task/:id/stop", taskController.Stop)
		api.PUT("/task/:id", taskController.Update)
		api.DELETE("/task/:id", taskController.Destroy)

		taskCrashController := new(TaskCrashController)
		api.GET("/crash/:id", taskCrashController.Download)

		storageItemController := new(StorageItemController)
		api.GET("/storage_item", middleware.Pagination, storageItemController.List)
		api.GET("/storage_item/:type", middleware.Pagination, storageItemController.ListByType)
		api.POST("/storage_item", storageItemController.Create)
		api.POST("/storage_item/exist", storageItemController.CreateExist)
		api.DELETE("/storage_item/:id", storageItemController.Destroy)

		api.GET("/task/:path1", TaskGetHandler)
		api.GET("/task/:path1/:path2", middleware.Pagination, TaskGetHandler)
	}
}

func prepareConfig() {
	viper.SetConfigType("YAML")
	data, err := ioutil.ReadFile("../config/config.yaml")

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Read 'config.yaml' fail: " + err.Error())
	}

	viper.ReadConfig(bytes.NewBuffer(data))
	viper.UnmarshalKey("server", config.ServerConf)
	config.ServerConf.CrashesPath = "./crashes"
	viper.UnmarshalKey("kubernetes", config.KubernetesConf)
	config.KubernetesConf.CheckTaskTime = 10
}

func TestMain(m *testing.M) {
	os.RemoveAll("./test.db")
	os.RemoveAll("./crashes")
	prepareConfig()
	prepareTestDB()
	prepareRouter()
	logger.Setup()
	service.Setup()
	m.Run()
	models.DB.Close()
}
