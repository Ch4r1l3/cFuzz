package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FuzzerController struct{}

func (fc *FuzzerController) List(c *gin.Context) {
	var fuzzers []models.Fuzzer
	if err := models.DB.Find(&fuzzers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}
	c.JSON(http.StatusOK, fuzzers)
}

func (fc *FuzzerController) Create(c *gin.Context) {
	var err error
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	name := c.PostForm("name")
	//check same name
	fuzzer, err := models.GetFuzzerByName(name)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fuzzer with same name exists"})
		return
	}

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fuzzer name cannot be empty"})
		return
	}
	isZipFile := false
	if strings.HasSuffix(header.Filename, ".zip") {
		isZipFile = true
	}
	tmpDir, err := ioutil.TempDir(config.ServerConf.TempPath, "fuzzer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error create temp directory"})
		return
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpDir)
		}
	}()
	var tempFile *os.File
	if isZipFile {
		tempFile, err = ioutil.TempFile(tmpDir, "fuzzer.*.zip")
	} else {
		tempFile, err = ioutil.TempFile(tmpDir, "fuzzer")
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error copy upload file"})
		return
	}
	fuzzer.Name = name
	fuzzer.Path = tempFile.Name()
	//unzip all file
	if isZipFile {
		err = utils.Unzip(tempFile.Name())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fuzzer.Path = filepath.Join(tmpDir, config.ServerConf.DefaultFuzzerName)
		if _, err = os.Stat(fuzzer.Path); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you should have fuzzer plugin in zip"})
			return
		}
	}
	err = os.Chmod(fuzzer.Path, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error change mode of fuzzer plugin"})
		return
	}

	models.DB.Save(&fuzzer)

	c.JSON(http.StatusOK, fuzzer)
}

func (fc *FuzzerController) Destroy(c *gin.Context) {
	name := c.Param("name")
	fuzzer, err := models.GetFuzzerByName(name)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(404)
		return
	}
	os.RemoveAll(filepath.Dir(fuzzer.Path))
	models.DB.Delete(fuzzer)
	c.String(http.StatusNoContent, "")
}
