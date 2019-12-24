package controller

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/bot/server/service"
	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
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
	models.DB.AutoMigrate(&models.Fuzzer{}, &models.Task{}, &models.TaskArgument{}, &models.TaskEnviroment{})

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
	prepareConfig()
	prepareTestDB()
	service.Setup()
	prepareRouter()
	m.Run()
	models.DB.Close()
	os.RemoveAll("./test.db")
}

func TestList(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/fuzzer").Expect().Status(http.StatusOK)
}

func TestFuzzer(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	err := ioutil.WriteFile("./fuzzer", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer")
	}()
	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("name", "afl").
		Expect().
		Status(http.StatusOK)
	e.GET("/fuzzer").Expect().
		Status(http.StatusOK).
		JSON().
		Array().First().Object().ValueEqual("name", "afl")
	e.DELETE("/fuzzer/afl").Expect().Status(http.StatusNoContent)

}

func TestFuzzerWithZipFile(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	err := ioutil.WriteFile(config.ServerConf.DefaultFuzzerName, []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(config.ServerConf.DefaultFuzzerName)
	}()
	file, err := os.Open(config.ServerConf.DefaultFuzzerName)
	if err != nil {
		t.Fatal(err)
	}
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		t.Fatal(err)
	}

	d, err := os.Create("./fuzzer.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer.zip")
	}()
	w := zip.NewWriter(d)
	writer, err := w.CreateHeader(header)
	if err != nil {
		w.Close()
		d.Close()
		t.Fatal(err)
	}
	_, err = io.Copy(writer, file)
	file.Close()
	w.Close()
	d.Close()
	if err != nil {
		t.Fatal(err)
	}

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer.zip").
		WithFormField("name", "aflzip").
		Expect().
		Status(http.StatusOK)

	e.GET("/fuzzer").Expect().
		Status(http.StatusOK).
		JSON().
		Array().First().Object().ValueEqual("name", "aflzip")
	e.DELETE("/fuzzer/aflzip").Expect().Status(http.StatusNoContent)

}

func TestFuzzerWithWrongZipFile1(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	filename := "abc"
	err := ioutil.WriteFile(filename, []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(filename)
	}()
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		t.Fatal(err)
	}

	d, err := os.Create("./fuzzer.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer.zip")
	}()
	w := zip.NewWriter(d)
	writer, err := w.CreateHeader(header)
	if err != nil {
		w.Close()
		d.Close()
		t.Fatal(err)
	}
	_, err = io.Copy(writer, file)
	file.Close()
	w.Close()
	d.Close()
	if err != nil {
		t.Fatal(err)
	}

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer.zip").
		WithFormField("name", "aflzip").
		Expect().
		Status(http.StatusBadRequest)

	e.GET("/fuzzer").Expect().
		Status(http.StatusOK).
		JSON().
		Array().Empty()

}

func TestFuzzerWithWrongZipFile2(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	filename := "abc.zip"
	err := ioutil.WriteFile(filename, []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(filename)
	}()

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", filename).
		WithFormField("name", "aflzip").
		Expect().
		Status(http.StatusBadRequest)

	e.GET("/fuzzer").Expect().
		Status(http.StatusOK).
		JSON().
		Array().Empty()

}
