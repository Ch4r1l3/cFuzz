package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type FuzzerUriReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

type FuzzerController struct{}

func (fc *FuzzerController) List(c *gin.Context) {
	fuzzers := []models.Fuzzer{}
	err := models.GetObjects(&fuzzers)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, fuzzers)
}

func (fc *FuzzerController) Create(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c)
		return
	}
	name := c.PostForm("name")
	if name == "" {
		utils.BadRequestWithMsg(c, "fuzzer name empty")
		return
	}
	if models.IsFuzzerExistsByName(name) {
		utils.BadRequestWithMsg(c, "fuzzer name exists")
		return
	}
	tempFile, err := ioutil.TempFile(config.ServerConf.TempPath, "fuzzer")
	if err != nil {
		utils.InternalErrorWithMsg(c, "error create temp file")
		return
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		tempFile.Close()
		os.RemoveAll(tempFile.Name())
		utils.InternalErrorWithMsg(c, "error copy upload file")
		return
	}
	tempFile.Close()
	fuzzer := models.Fuzzer{
		Name: name,
		Path: tempFile.Name(),
	}
	err = models.InsertObject(&fuzzer)
	if err != nil {
		os.RemoveAll(tempFile.Name())
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, fuzzer)
}

func (fc *FuzzerController) Destroy(c *gin.Context) {
	var req FuzzerUriReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	err = models.DeleteObjectByID(models.Fuzzer{}, req.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}
