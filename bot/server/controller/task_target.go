package controller

import (
	"errors"
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

type TaskTargetController struct{}

func (ttc *TaskTargetController) Create(c *gin.Context) {
	var Err error
	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exist, please create one")
		return
	}
	if task.TargetDir != "" {
		if _, err = os.Stat(task.TargetDir); !os.IsNotExist(err) {
			os.RemoveAll(task.TargetDir)
		}
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequestWithMsg(c, "please upload file")
		return
	}
	isZipFile := false
	targetRelPath := c.PostForm("targetRelPath")
	if strings.HasSuffix(header.Filename, ".zip") {
		isZipFile = true
		if targetRelPath == "" {
			utils.BadRequestWithMsg(c, "lack of targetRelPath if you post zip file")
			return
		}
	}
	tmpDir, err := ioutil.TempDir(config.ServerConf.TempPath, "target")
	if err != nil {
		utils.InternalErrorWithMsg(c, "error create temp directory")
		return
	}
	defer func() {
		if Err != nil {
			os.RemoveAll(tmpDir)
			models.DB.Model(task).Update("TargetDir", "")
			models.DB.Model(task).Update("TargetPath", "")
		}
	}()

	var tempFile *os.File
	if isZipFile {
		tempFile, err = ioutil.TempFile(tmpDir, "target.*.zip")
		if err != nil {
			Err = err
			utils.InternalErrorWithMsg(c, "error create temp file")
			return
		}
	} else {
		tempFile, err = ioutil.TempFile(tmpDir, "target")
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
	targetPath := tempFile.Name()
	if isZipFile {
		err = utils.Unzip(tempFile.Name())
		if err != nil {
			Err = err
			utils.BadRequest(c)
			return
		}
		targetPath = filepath.Join(tmpDir, targetRelPath)
		targetPath = filepath.Clean(targetPath)
		relPath, err := filepath.Rel(tmpDir, targetPath)
		if err != nil {
			Err = err
			utils.BadRequestWithMsg(c, "error get rel of targetRelPath")
			return
		}
		if strings.Contains(relPath, "..") {
			Err = errors.New("relPath wrong")
			utils.BadRequestWithMsg(c, "path can only under this temp directory")
			return
		}
		if _, err = os.Stat(targetPath); os.IsNotExist(err) {

			Err = err
			utils.BadRequestWithMsg(c, "target not exists in zip file")

			return
		}
		os.RemoveAll(tempFile.Name())
	}
	err = os.Chmod(targetPath, 0755)
	if err != nil {
		Err = err
		utils.InternalErrorWithMsg(c, "error change mode of fuzzer plugin")
		return
	}

	models.DB.Model(task).Update("TargetPath", targetPath)
	models.DB.Model(task).Update("TargetDir", tmpDir)
	c.JSON(http.StatusOK, "")
}

func (ttc *TaskTargetController) Destroy(c *gin.Context) {
	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exist, please create one")
		return
	}

	if task.TargetDir != "" {
		if _, err = os.Stat(task.TargetDir); !os.IsNotExist(err) {
			os.RemoveAll(task.TargetDir)
		}

		models.DB.Model(task).Update("TargetPath", "")
		models.DB.Model(task).Update("TargetDir", "")
	}
	c.JSON(http.StatusNoContent, "")

}
