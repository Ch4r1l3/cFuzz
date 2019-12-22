package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type FuzzerController struct{}

func (fc *FuzzerController) List(c *gin.Context) {
	var fuzzers []models.Fuzzer
	if err := models.DB.Find(&fuzzers).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, fuzzers)
}

func (fc *FuzzerController) Create(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}
	name := c.PostForm("name")
	if name == "" {
		c.String(http.StatusBadRequest, "fuzzer name cannot be empty")
		return
	}
	isZipFile := false
	if strings.HasSuffix(header.Filename, ".zip") {
		isZipFile = true
	}
	tmpDir, err := ioutil.TempDir(config.ServerConf.FuzzerStorePath, "fuzzer")
	if err != nil {
		c.String(http.StatusInternalServerError, "error create temp directory")
		return
	}
	var tempFile *os.File
	if isZipFile {
		tempFile, err = ioutil.TempFile(tmpDir, "fuzzer.*.zip")
	} else {
		tempFile, err = ioutil.TempFile(tmpDir, "fuzzer")
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "error copy upload file")
		return
	}
	fuzzer := models.Fuzzer{
		Name: name,
		Path: tempFile.Name(),
	}
	models.DB.Save(&fuzzer)

	c.JSON(200, fuzzer)
}

func (fc *FuzzerController) Destroy(c *gin.Context) {
	name := c.Param("name")
	var fuzzer models.Fuzzer
	if err := models.DB.Where("name = ?", name).First(&fuzzer).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatus(404)
		return
	}
	models.DB.Delete(&fuzzer)
	c.String(http.StatusNoContent, "")

}
