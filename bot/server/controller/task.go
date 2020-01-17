package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/bot/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type TaskController struct{}

type TaskCreateReq struct {
	FuzzerID      uint64            `json:"fuzzerID" binding:"required"`
	MaxTime       int               `json:"maxTime" binding:"required"`
	FuzzCycleTime uint64            `json:"fuzzCycleTime" binding:"required"`
	Arguments     map[string]string `json:"arguments"`
	Environments  []string          `json:"environments"`
}

type TaskUpdateReq struct {
	Status string `json:"status" binding:"required"`
}

type TaskIDReq struct {
	ID uint64 `uri:"id" binding:"required"`
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
		"corpusDir":     task.CorpusDir,
		"targetDir":     task.TargetDir,
		"targetPath":    task.TargetPath,
		"fuzzerID":      task.FuzzerID,
		"maxTime":       task.MaxTime,
		"fuzzCycleTime": task.FuzzCycleTime,
		"status":        task.Status,
		"arguments":     arguments,
		"environments":  environments,
	})

}

func (tc *TaskController) Create(c *gin.Context) {
	//var task models.Task
	var req TaskCreateReq
	err := c.BindJSON(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
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
	task.FuzzCycleTime = req.FuzzCycleTime
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
		service.Fuzz(fuzzer.Path, task.TargetPath, task.CorpusDir, task.MaxTime, int(task.FuzzCycleTime), arguments, environments)

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
