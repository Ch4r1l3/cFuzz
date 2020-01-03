package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
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
	var err error
	name := c.PostForm("name")
	if name == "" {
		utils.BadRequestWithMsg(c, "fuzzer name empty")
		return
	}
	if models.IsFuzzerExistsByName(name) {
		utils.BadRequestWithMsg(c, "fuzzer name exists")
		return
	}
	var tempFile string
	if tempFile, err = utils.SaveTempFile(c, "file", "fuzzer"); err != nil {
		return
	}
	fuzzer := models.Fuzzer{
		Name: name,
		Path: tempFile,
	}
	err = models.InsertObject(&fuzzer)
	if err != nil {
		os.RemoveAll(tempFile)
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
	if !models.IsObjectExistsByID(&models.Fuzzer{}, req.ID) {
		utils.NotFound(c)
		return
	}
	var fuzzer models.Fuzzer
	if err = models.GetObjectByID(&fuzzer, req.ID); err != nil {
		utils.DBError(c)
		return
	}
	os.RemoveAll(fuzzer.Path)
	err = models.DeleteObjectByID(models.Fuzzer{}, req.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}
