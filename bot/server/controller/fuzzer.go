package controller

import (
	"archive/zip"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
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
	var fuzzer models.Fuzzer
	if err := models.DB.Where("name = ?", name).First(&fuzzer).Error; err == nil {
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
	tmpDir, err := ioutil.TempDir(config.ServerConf.FuzzerStorePath, "fuzzer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error create temp directory"})
		return
	}
	go func() {
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
		reader, err := zip.OpenReader(tempFile.Name())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error open zip file"})
			return
		}
		for _, file := range reader.File {
			tmpPath := filepath.Join(tmpDir, file.Name)
			tmpPath = filepath.Clean(tmpPath)
			relPath, err := filepath.Rel(tmpDir, tmpPath)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "error get rel of zip file"})
				return
			}
			if filepath.Join(tmpDir, relPath) != tmpPath {
				c.JSON(http.StatusBadRequest, gin.H{"error": "zipslip not work here"})
			}
			if file.FileInfo().IsDir() {
				os.MkdirAll(tmpPath, file.Mode())
				continue
			}

			fileReader, err := file.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "error open " + file.Name + " of the zip file"})
				return
			}
			defer fileReader.Close()

			targetFile, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error create file when unzip"})
				return
			}
			defer targetFile.Close()

			if _, err = io.Copy(targetFile, fileReader); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error copy file"})
				return
			}
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
	os.RemoveAll(filepath.Dir(fuzzer.Path))
	models.DB.Delete(&fuzzer)
	c.String(http.StatusNoContent, "")
}
