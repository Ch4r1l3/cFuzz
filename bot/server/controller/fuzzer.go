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
	"strconv"
	"strings"
)

type FuzzerController struct{}

func (fc *FuzzerController) List(c *gin.Context) {
	var fuzzers []models.Fuzzer
	if err := models.DB.Find(&fuzzers).Error; err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, fuzzers)
}

func (fc *FuzzerController) Create(c *gin.Context) {
	var Err error
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c)
		return
	}
	name := c.PostForm("name")
	if name == "" {
		utils.BadRequestWithMsg(c, "fuzzer name cannot be empty")
		return
	}
	isZipFile := false
	if strings.HasSuffix(header.Filename, ".zip") {
		isZipFile = true
	}
	tmpDir, err := ioutil.TempDir(config.ServerConf.TempPath, "fuzzer")
	if err != nil {
		utils.InternalErrorWithMsg(c, "error create temp directory")
		return
	}
	defer func() {
		if Err != nil {
			os.RemoveAll(tmpDir)
		}
	}()
	var tempFile *os.File
	if isZipFile {
		tempFile, err = ioutil.TempFile(tmpDir, "fuzzer.*.zip")
		if err != nil {
			Err = err
			utils.InternalErrorWithMsg(c, "error create temp file")
			return
		}
	} else {
		tempFile, err = ioutil.TempFile(tmpDir, "fuzzer")
		if err != nil {
			Err = err
			utils.InternalErrorWithMsg(c, "error create temp file")
			return
		}
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		Err = err
		tempFile.Close()
		utils.InternalErrorWithMsg(c, "error copy upload file")
		return
	}
	tempFile.Close()
	fuzzer := models.Fuzzer{
		Name: name,
		Path: tempFile.Name(),
	}
	//unzip all file
	if isZipFile {
		err = utils.Unzip(tempFile.Name())
		if err != nil {
			Err = err
			utils.BadRequestWithMsg(c, err.Error())
			return
		}
		fuzzer.Path = filepath.Join(tmpDir, config.ServerConf.DefaultFuzzerName)
		if _, err = os.Stat(fuzzer.Path); os.IsNotExist(err) {
			Err = err
			utils.BadRequestWithMsg(c, "you should have fuzzer plugin in zip")
			return
		}
		os.RemoveAll(tempFile.Name())
	}
	err = os.Chmod(fuzzer.Path, 0755)
	if err != nil {
		Err = err
		utils.InternalErrorWithMsg(c, "error change mode of fuzzer plugin")
		return
	}

	models.DB.Save(&fuzzer)

	c.JSON(http.StatusOK, fuzzer)
}

func (fc *FuzzerController) Destroy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	fuzzer, err := models.GetFuzzerByID(uint64(id))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(404)
		return
	}
	os.RemoveAll(filepath.Dir(fuzzer.Path))
	models.DB.Delete(fuzzer)
	c.String(http.StatusNoContent, "")
}
