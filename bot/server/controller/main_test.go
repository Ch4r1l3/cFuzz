package controller

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
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
	models.DB.AutoMigrate(&models.Fuzzer{}, &models.Task{}, &models.TaskArgument{}, &models.TaskEnvironment{})

}

func prepareRouter() {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(config.ServerConf.RunMode)
	r.MaxMultipartMemory = config.ServerConf.UploadFileLimit << 20

	fuzzerController := new(FuzzerController)
	r.GET("/fuzzer", fuzzerController.List)
	r.POST("/fuzzer", fuzzerController.Create)
	r.DELETE("/fuzzer/:name", fuzzerController.Destroy)

	taskController := new(TaskController)
	r.GET("/task", taskController.Retrieve)
	r.POST("/task", taskController.Create)
	r.PUT("/task", taskController.Update)
	r.DELETE("/task", taskController.Destroy)

	taskCrashController := new(TaskCrashController)
	r.GET("/task/crash", taskCrashController.List)

	taskCorpusController := new(TaskCorpusController)
	r.POST("/task/corpus", taskCorpusController.Create)
	r.DELETE("/task/corpus", taskCorpusController.Destroy)

	taskTargetController := new(TaskTargetController)
	r.POST("/task/target", taskTargetController.Create)
	r.DELETE("/task/target", taskTargetController.Destroy)
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
	service.Setup()
	prepareRouter()
	m.Run()
	models.DB.Close()
	os.RemoveAll("./test.db")
}
