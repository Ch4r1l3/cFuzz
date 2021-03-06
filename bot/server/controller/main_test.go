package controller

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/logger"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/bot/server/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	"io"
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
	models.DB.AutoMigrate(&models.StorageItem{}, &models.Task{}, &models.TaskCrash{}, &models.TaskArgument{}, &models.TaskEnvironment{}, &models.TaskFuzzResult{}, &models.TaskFuzzResultStat{})

}

func prepareRouter() {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(config.ServerConf.RunMode)
	r.MaxMultipartMemory = config.ServerConf.UploadFileLimit << 20

	storageItemController := new(StorageItemController)
	r.GET("/storage_item", storageItemController.List)
	r.GET("/storage_item/:id", storageItemController.Retrieve)
	r.POST("/storage_item", storageItemController.Create)
	r.POST("/storage_item/exist", storageItemController.CreateExist)
	r.DELETE("/storage_item/:id", storageItemController.Destroy)

	taskController := new(TaskController)
	r.GET("/task", taskController.Retrieve)
	r.POST("/task", taskController.Create)
	r.POST("/task/start", taskController.StartFuzz)
	r.POST("/task/stop", taskController.StopFuzz)
	r.DELETE("/task", taskController.Destroy)

	taskCrashController := new(TaskCrashController)
	r.GET("/task/crash", taskCrashController.List)
	r.GET("/task/crash/:id", taskCrashController.Download)

	taskResultController := new(TaskResultController)
	r.GET("/task/result", taskResultController.Retrieve)
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

func createZipfile(zipfileName string, fileName string, fileContent []byte) error {
	err := ioutil.WriteFile(fileName, fileContent, 0755)
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(fileName)
	}()
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	d, err := os.Create(zipfileName)
	if err != nil {
		return err
	}

	w := zip.NewWriter(d)
	writer, err := w.CreateHeader(header)
	if err != nil {
		w.Close()
		d.Close()
		os.RemoveAll(zipfileName)
		return err
	}
	_, err = io.Copy(writer, file)
	file.Close()
	w.Close()
	d.Close()
	if err != nil {
		os.RemoveAll(zipfileName)
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	prepareConfig()
	prepareTestDB()
	logger.Setup()
	service.Setup()
	prepareRouter()
	m.Run()
	models.DB.Close()
	os.RemoveAll("./test.db")
}
