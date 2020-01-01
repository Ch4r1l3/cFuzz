package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type FuzzerUriReq struct {
	Id uint64 `uri:"id" binding:"required"`
}

type FuzzerController struct{}

func (fc *FuzzerController) List(c *gin.Context) {
	fuzzers := []models.Fuzzer{}
	err := models.GetObjects(&fuzzers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}
	c.JSON(http.StatusOK, fuzzers)
}

func (fc *FuzzerController) Create(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fuzzer name empty"})
		return
	}
	if models.IsFuzzerExists(name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fuzzer name exists"})
		return
	}
	tempFile, err := ioutil.TempFile(config.ServerConf.TempPath, "fuzzer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error create temp file"})
		return
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		tempFile.Close()
		os.RemoveAll(tempFile.Name())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error copy upload file"})
		return
	}
	tempFile.Close()
	err = models.InsertObject(&models.Fuzzer{
		Name: name,
		Path: tempFile.Name(),
	})
	if err != nil {
		os.RemoveAll(tempFile.Name())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, "")
}

func (fc *FuzzerController) Destroy(c *gin.Context) {
	var req FuzzerUriReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	err = models.DeleteObjectById(models.Fuzzer{}, req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}
	c.JSON(http.StatusNoContent, "")
}
