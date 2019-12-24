package controller

import (
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
	FuzzerName   string            `json:"fuzzerName"`
	MaxTime      int               `json:"maxTime"`
	Arguments    map[string]string `json:"arguments"`
	Environments []string          `json:"environments"`
}

type TaskUpdateReq struct {
	Status string `json:"status"`
}

func (tc *TaskController) Retrieve(c *gin.Context) {
	task, err1 := models.GetTask()
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "create task first"})
		return
	}
	arguments, err2 := models.GetArguments()
	environments, err3 := models.GetEnvironments()
	if err2 != nil || err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"corpusDir":    task.CorpusDir,
		"targetDir":    task.TargetDir,
		"targetPath":   task.TargetPath,
		"fuzzerName":   task.FuzzerName,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	task, err := models.GetTask()
	if err == nil && task.Status == config.TASK_RUNNING {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task running"})
		return
	}
	_, err = models.GetFuzzerByName(req.FuzzerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fuzzer not exists"})
		return
	}

	if req.MaxTime <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fuzz run time should longger than 0s"})
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
	task.Status = config.TASK_CREATED
	task.FuzzerName = req.FuzzerName
	task.MaxTime = req.MaxTime
	models.DB.Create(&task)

	//create arguments
	err = models.InsertArguments(req.Arguments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	//create environments
	err = models.InsertEnvironments(req.Environments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (tc *TaskController) Update(c *gin.Context) {
	var req TaskUpdateReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	task, err := models.GetTask()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not exists"})
		return
	}
	if task.Status == config.TASK_RUNNING && req.Status == config.TASK_STOPPED {
		service.StopFuzz()
	} else if task.Status == config.TASK_CREATED && req.Status == config.TASK_RUNNING {
		//check plugin and target
		if _, err = os.Stat(task.CorpusDir); task.CorpusDir == "" || os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "you should upload corpus"})
			return
		}
		if _, err = os.Stat(task.TargetPath); task.TargetPath == "" || os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "you should upload target"})
			return
		}
		fuzzer, err := models.GetFuzzerByName(task.FuzzerName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fuzzer not exists"})
			return
		}
		arguments, err := models.GetArguments()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		environments, err := models.GetEnvironments()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		service.Fuzz(fuzzer.Path, task.TargetPath, task.CorpusDir, task.MaxTime, config.ServerConf.DefaultFuzzTime, arguments, environments)

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong status"})
		return
	}

	c.JSON(http.StatusOK, "")

}

func (tc *TaskController) Destroy(c *gin.Context) {
	service.StopFuzz()
	models.DB.Delete(&models.Task{})
	models.DB.Delete(&models.TaskArgument{})
	models.DB.Delete(&models.TaskEnvironment{})
	c.JSON(http.StatusNoContent, "")
}

type TaskCrashController struct{}

func (tcc *TaskCrashController) List(c *gin.Context) {
	crashes, err := models.GetCrashes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, crashes)
}

type TaskCorpusController struct{}

func (t *TaskCorpusController) Create(c *gin.Context) {
	var err error
	task, err := models.GetTask()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not exist, please create one"})
		return
	}
	if task.CorpusDir != "" {
		if _, err = os.Stat(task.CorpusDir); !os.IsNotExist(err) {
			os.RemoveAll(task.CorpusDir)
		}
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please upload file"})
		return
	}
	isZipFile := false
	if strings.HasSuffix(header.Filename, ".zip") {
		isZipFile = true
	}
	tmpDir, err := ioutil.TempDir(config.ServerConf.TempPath, "corpus")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error create temp directory"})
		return
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpDir)
			models.DB.Model(task).Update("CorpusDir", "")
		}
	}()
	var tempFile *os.File
	if isZipFile {
		tempFile, err = ioutil.TempFile(tmpDir, "corpus.*.zip")
	} else {
		tempFile, err = ioutil.TempFile(tmpDir, "corpus")
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error copy upload file"})
		return
	}
	if isZipFile {
		err = utils.Unzip(tempFile.Name())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	models.DB.Model(task).Update("CorpusDir", tmpDir)
	c.JSON(http.StatusOK, "")

}

func (t *TaskCorpusController) Destroy(c *gin.Context) {
	task, err := models.GetTask()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not exist, please create one"})
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
	task, err := models.GetTask()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not exist, please create one"})
		return
	}
	if task.TargetPath != "" {
		if _, err = os.Stat(task.TargetPath); !os.IsNotExist(err) {
			os.RemoveAll(filepath.Dir(task.TargetPath))
		}
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please upload file"})
		return
	}
	isZipFile := false
	targetRelPath := c.PostForm("targetRelPath")
	if strings.HasSuffix(header.Filename, ".zip") {
		isZipFile = true
		if targetRelPath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lack of targetRelPath if you post zip file"})
			return
		}
	}
	tmpDir, err := ioutil.TempDir(config.ServerConf.TempPath, "target")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error create temp directory"})
		return
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpDir)
			models.DB.Model(task).Update("TargetDir", "")
			models.DB.Model(task).Update("TargetPath", "")
		}
	}()
	var tempFile *os.File
	if isZipFile {
		tempFile, err = ioutil.TempFile(tmpDir, "target.*.zip")
	} else {
		tempFile, err = ioutil.TempFile(tmpDir, "target")
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error copy upload file"})
		return
	}
	targetPath := tempFile.Name()
	if isZipFile {
		err = utils.Unzip(tempFile.Name())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		targetPath = filepath.Join(tmpDir, targetRelPath)
		targetPath = filepath.Clean(targetPath)
		relPath, err := filepath.Rel(tmpDir, targetPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error get rel of targetRelPath"})
			return
		}
		if filepath.Join(tmpDir, relPath) != targetPath {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path can only under this temp directory"})
			return
		}
		if _, err = os.Stat(targetPath); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target not exists in zip file"})
			return
		}

	}

	models.DB.Model(task).Update("TargetPath", targetPath)
	models.DB.Model(task).Update("TargetDir", tmpDir)
	c.JSON(http.StatusOK, "")

}

func (ttc *TaskTargetController) Destroy(c *gin.Context) {
	task, err := models.GetTask()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not exist, please create one"})
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
