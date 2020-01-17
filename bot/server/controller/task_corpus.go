package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type TaskCorpusController struct{}

func (t *TaskCorpusController) Create(c *gin.Context) {
	var Err error
	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exist, please create one")
		return
	}
	if task.CorpusDir != "" {
		if _, err = os.Stat(task.CorpusDir); !os.IsNotExist(err) {
			os.RemoveAll(task.CorpusDir)
		}
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequestWithMsg(c, "please upload file")
		return
	}
	isZipFile := false
	if strings.HasSuffix(header.Filename, ".zip") {
		isZipFile = true
	}
	tmpDir, err := ioutil.TempDir(config.ServerConf.TempPath, "corpus")
	if err != nil {
		utils.InternalErrorWithMsg(c, "error create temp directory")
		return
	}
	defer func() {
		if Err != nil {
			os.RemoveAll(tmpDir)
			models.DB.Model(task).Update("CorpusDir", "")
		}
	}()
	var tempFile *os.File
	if isZipFile {
		tempFile, err = ioutil.TempFile(tmpDir, "corpus.*.zip")
		if err != nil {
			Err = err
			utils.InternalErrorWithMsg(c, "error create temp file")
			return
		}
	} else {
		tempFile, err = ioutil.TempFile(tmpDir, "corpus")
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
	if isZipFile {
		err = utils.Unzip(tempFile.Name())
		if err != nil {
			Err = err
			utils.BadRequestWithMsg(c, err.Error())
			return
		}
		os.RemoveAll(tempFile.Name())
	}

	models.DB.Model(task).Update("CorpusDir", tmpDir)
	c.JSON(http.StatusOK, "")
}

func (t *TaskCorpusController) Destroy(c *gin.Context) {
	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exist, please create one")
		return
	}
	if task.CorpusDir != "" {
		if _, err = os.Stat(task.CorpusDir); !os.IsNotExist(err) {
			os.RemoveAll(task.CorpusDir)
		}
		models.DB.Model(task).Update("CorpusDir", "")
	}
	c.JSON(http.StatusNoContent, "")

}
