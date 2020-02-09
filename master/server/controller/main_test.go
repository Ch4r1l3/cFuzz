package controller

import (
	"bytes"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/middleware"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	r *gin.Engine
	s *httptest.Server
)

func prepareTestDB() {
	var err error
	models.DB, err = gorm.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}

	models.DB.SingularTable(true)
	models.DB.AutoMigrate(&models.Image{}, &models.Task{}, &models.StorageItem{}, &models.TaskEnvironment{}, &models.TaskArgument{}, &models.TaskCrash{}, &models.TaskFuzzResult{}, &models.TaskFuzzResultStat{}, &models.User{})

}

func prepareRouter() {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode("debug")

	userController := new(UserController)

	api := r.Group("/api")
	api.Use(middleware.Auth)
	{
		imageController := new(ImageController)
		api.GET("/image", middleware.Pagination, imageController.List)
		api.GET("/image/:id", imageController.Retrieve)
		api.POST("/image", middleware.CheckUserExist, imageController.Create)
		api.PUT("/image/:id", middleware.CheckUserExist, imageController.Update)
		api.DELETE("/image/:id", middleware.CheckUserExist, imageController.Destroy)

		taskController := new(TaskController)
		api.GET("/task", middleware.Pagination, taskController.List)
		api.POST("/task", middleware.CheckUserExist, taskController.Create)
		api.POST("/task/:id/start", middleware.CheckUserExist, taskController.Start)
		api.POST("/task/:id/stop", middleware.CheckUserExist, taskController.Stop)
		api.PUT("/task/:id", middleware.CheckUserExist, taskController.Update)
		api.DELETE("/task/:id", middleware.CheckUserExist, taskController.Destroy)

		taskCrashController := new(TaskCrashController)
		api.GET("/crash/:id", taskCrashController.Download)

		storageItemController := new(StorageItemController)
		api.GET("/storage_item", middleware.Pagination, storageItemController.List)
		api.GET("/storage_item/:type", middleware.Pagination, storageItemController.ListByType)
		api.POST("/storage_item", middleware.CheckUserExist, storageItemController.Create)
		api.POST("/storage_item/exist", middleware.CheckUserExist, storageItemController.CreateExist)
		api.DELETE("/storage_item/:id", middleware.CheckUserExist, storageItemController.Destroy)

		api.GET("/task/:path1", TaskGetHandler)
		api.GET("/task/:path1/:path2", middleware.Pagination, TaskGetHandler)

		api.GET("/user/info", userController.Info)
		api.GET("/user", middleware.AdminOnly, middleware.Pagination, userController.List)
		api.POST("/user", middleware.AdminOnly, userController.Create)
		api.PUT("/user/:id", userController.Update)
		api.DELETE("/user/:id", middleware.AdminOnly, userController.Delete)
	}
	r.POST("/api/user/login", userController.Login)
	r.GET("/api/user/logout", userController.Logout)
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
	config.ServerConf.SigningKey = "cfuzz"
}

func getAdminExpect(t *testing.T) *httpexpect.Expect {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  s.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Jar: httpexpect.NewJar(),
		},
	})
	loginData := map[string]interface{}{
		"username": "admin",
		"password": "123456",
	}
	e.POST("/api/user/login").WithJSON(loginData).Expect().Status(http.StatusOK)
	return e
}

func getExpect(t *testing.T) *httpexpect.Expect {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  s.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Jar: httpexpect.NewJar(),
		},
	})
	loginData := map[string]interface{}{
		"username": "abc",
		"password": "123456",
	}
	e.POST("/api/user/login").WithJSON(loginData).Expect().Status(http.StatusOK)
	return e
}

func TestMain(m *testing.M) {
	os.RemoveAll("./test.db")
	os.RemoveAll("./crashes")
	prepareConfig()
	prepareTestDB()
	prepareRouter()
	logger.Setup()
	service.Setup()
	s = httptest.NewServer(r)
	defer s.Close()
	if err := models.CreateUser("admin", "123456", true); err != nil {
		log.Fatal("create user error")
	}
	if err := models.CreateUser("abc", "123456", false); err != nil {
		log.Fatal("create user error")
	}
	m.Run()
	models.DB.Close()
}
