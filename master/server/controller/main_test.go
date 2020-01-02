package controller

import (
	"bytes"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
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
	models.DB.AutoMigrate(&models.Dockerfile{}, &models.Task{}, &models.TaskTarget{}, &models.TaskCorpus{}, &models.TaskEnvironment{}, &models.TaskArgument{}, &models.TaskCrash{}, &models.TaskFuzzResult{}, &models.TaskFuzzResultStat{}, &models.Fuzzer{})

}

func prepareRouter() {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode("debug")

	dockerfileController := new(DockerfileController)
	r.GET("/dockerfile", dockerfileController.List)
	r.POST("/dockerfile", dockerfileController.Create)
	r.PUT("/dockerfile/:id", dockerfileController.Update)
	r.DELETE("/dockerfile/:id", dockerfileController.Destroy)

	taskController := new(TaskController)
	r.GET("/task", taskController.List)
	r.POST("/task", taskController.Create)
	r.PUT("/task/:id", taskController.Update)
	r.DELETE("/task/:path1", TaskDeleteHandler)

	taskCorpusController := new(TaskCorpusController)
	r.GET("/task/:taskid/corpus", taskCorpusController.List)
	r.POST("/task/:taskid/corpus", taskCorpusController.Create)
	r.DELETE("/task/:path1/:path2/:path3", TaskDeleteHandler)

	taskTargetController := new(TaskTargetController)
	r.GET("/task/:taskid/target", taskTargetController.Retrieve)
	r.POST("/task/:taskid/target", taskTargetController.Create)
	r.DELETE("/task/:path1/:path2", TaskDeleteHandler)

	fuzzerController := new(FuzzerController)
	r.GET("/fuzzer", fuzzerController.List)
	r.POST("/fuzzer", fuzzerController.Create)
	r.DELETE("/fuzzer/:id", fuzzerController.Destroy)
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
}

func TestMain(m *testing.M) {
	os.RemoveAll("./test.db")
	prepareConfig()
	prepareTestDB()
	prepareRouter()
	m.Run()
	models.DB.Close()
}
