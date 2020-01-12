package controller

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/bot/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type TaskController struct{}

type TaskCreateReq struct {
	FuzzerID     uint64            `json:"fuzzerID" binding:"required"`
	MaxTime      int               `json:"maxTime" binding:"required"`
	Arguments    map[string]string `json:"arguments"`
	Environments []string          `json:"environments"`
}

type TaskUpdateReq struct {
	Status string `json:"status" binding:"required"`
}

func (tc *TaskController) Retrieve(c *gin.Context) {
	task, err1 := models.GetTask()
	if err1 != nil {
		utils.BadRequestWithMsg(c, "create task first")
		return
	}
	arguments, err2 := models.GetArguments()
	environments, err3 := models.GetEnvironments()
	if err2 != nil || err3 != nil {
		utils.DBError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"corpusDir":    task.CorpusDir,
		"targetDir":    task.TargetDir,
		"targetPath":   task.TargetPath,
		"fuzzerID":     task.FuzzerID,
		"maxTime":      task.MaxTime,
		"status":       task.Status,
		"arguments":    arguments,
		"environments": environments,
	})

}

func (tc *TaskController) Create(c *gin.Context) {
	//var task models.Task
	var req TaskCreateReq
	err := c.BindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	task, err := models.GetTask()
	if err == nil && task.Status == models.TASK_RUNNING {
		utils.BadRequestWithMsg(c, "task running")
		return
	}
	_, err = models.GetFuzzerByID(req.FuzzerID)
	if err != nil {
		utils.BadRequestWithMsg(c, "fuzzer not exists")
		return
	}

	if req.MaxTime <= 0 {
		utils.BadRequestWithMsg(c, "fuzz run time should longger than 0s")
		return
	}

	//remove corpus dir and target path
	if _, err = os.Stat(task.CorpusDir); task.CorpusDir != "" && !os.IsNotExist(err) {
		os.RemoveAll(task.CorpusDir)
	}

	if _, err = os.Stat(task.TargetDir); task.TargetDir != "" && !os.IsNotExist(err) {
		os.RemoveAll(task.TargetDir)
	}
	//clear tasks and others
	models.DB.Delete(&task)
	models.DB.Delete(&models.TaskArgument{})
	models.DB.Delete(&models.TaskEnvironment{})

	//create task
	task.CorpusDir = ""
	task.TargetPath = ""
	task.TargetDir = ""
	task.Status = models.TASK_CREATED
	task.FuzzerID = req.FuzzerID
	task.MaxTime = req.MaxTime
	models.DB.Create(&task)

	//create arguments
	err = models.InsertArguments(req.Arguments)
	if err != nil {
		utils.DBError(c)
		return
	}

	//create environments
	err = models.InsertEnvironments(req.Environments)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (tc *TaskController) Update(c *gin.Context) {
	var req TaskUpdateReq
	err := c.BindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exists")
		return
	}
	if task.Status == models.TASK_RUNNING && req.Status == models.TASK_STOPPED {
		service.StopFuzz()
		models.DB.Model(task).Update("Status", models.TASK_STOPPED)

	} else if task.Status == models.TASK_CREATED && req.Status == models.TASK_RUNNING {
		//check plugin and target
		if _, err = os.Stat(task.CorpusDir); task.CorpusDir == "" || os.IsNotExist(err) {
			utils.BadRequestWithMsg(c, "you should upload corpus")
			return
		}
		if _, err = os.Stat(task.TargetPath); task.TargetPath == "" || os.IsNotExist(err) {
			utils.BadRequestWithMsg(c, "you should upload target")
			return
		}
		fuzzer, err := models.GetFuzzerByID(task.FuzzerID)
		if err != nil {
			utils.BadRequestWithMsg(c, "fuzzer not exists")
			return
		}
		arguments, err := models.GetArguments()
		if err != nil {
			utils.DBError(c)
			return
		}
		environments, err := models.GetEnvironments()
		if err != nil {
			utils.DBError(c)
			return
		}
		service.Fuzz(fuzzer.Path, task.TargetPath, task.CorpusDir, task.MaxTime, config.ServerConf.DefaultFuzzTime, arguments, environments)

		models.DB.Model(task).Update("Status", models.TASK_RUNNING)

	} else {
		utils.BadRequestWithMsg(c, "wrong status")
		return
	}

	c.JSON(http.StatusOK, "")

}

func (tc *TaskController) Destroy(c *gin.Context) {
	service.StopFuzz()
	task, err := models.GetTask()
	if err != nil {
		utils.DBError(c)
		return
	}
	if task.CorpusDir != "" {
		if _, err = os.Stat(task.CorpusDir); !os.IsNotExist(err) {
			os.RemoveAll(task.CorpusDir)
		}
	}
	if task.TargetDir != "" {
		if _, err = os.Stat(task.TargetDir); !os.IsNotExist(err) {
			os.RemoveAll(task.TargetDir)
		}
	}
	models.DB.Delete(&models.Task{})
	models.DB.Delete(&models.TaskArgument{})
	models.DB.Delete(&models.TaskEnvironment{})
	c.JSON(http.StatusNoContent, "")
}

type TaskCrashController struct{}

func (tcc *TaskCrashController) List(c *gin.Context) {
	crashes, err := models.GetCrashes()
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, crashes)
}

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

type TaskResultController struct{}

func (trc *TaskResultController) Retrieve(c *gin.Context) {
	result, stats, err := models.GetFuzzResult()
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"command":      result.Command,
		"timeExecuted": result.TimeExecuted,
		"stats":        stats,
	})
}
