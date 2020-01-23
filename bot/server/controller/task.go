package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/bot/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskController struct{}

type TaskCreateReq struct {
	FuzzerID      uint64            `json:"fuzzerID" binding:"required"`
	CorpusID      uint64            `json:"corpusID" binding:"required"`
	TargetID      uint64            `json:"targetID" binding:"required"`
	MaxTime       int               `json:"maxTime" binding:"required"`
	FuzzCycleTime uint64            `json:"fuzzCycleTime" binding:"required"`
	Arguments     map[string]string `json:"arguments"`
	Environments  []string          `json:"environments"`
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
		"fuzzerID":      task.FuzzerID,
		"corpusID":      task.CorpusID,
		"targetID":      task.TargetID,
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

	ids := []uint64{req.FuzzerID, req.CorpusID, req.TargetID}
	types := []string{models.Fuzzer, models.Corpus, models.Target}
	for i, _ := range ids {
		ok, err := models.IsStorageItemExistByID(ids[i])
		if err != nil {
			utils.DBError(c)
			return
		}
		if !ok {
			utils.BadRequestWithMsg(c, types[i]+" not exists")
			return
		}
		storageItem, err := models.GetStorageItemByID(ids[i])
		if err != nil {
			utils.DBError(c)
			return
		}
		if storageItem.Type != types[i] {
			utils.BadRequestWithMsg(c, "type wrong")
			return
		}
	}

	if req.MaxTime <= 0 {
		utils.BadRequestWithMsg(c, "fuzz run time should longger than 0s")
		return
	}

	//clear tasks and others
	models.DB.Delete(&task)
	models.DB.Delete(&models.TaskArgument{})
	models.DB.Delete(&models.TaskEnvironment{})

	//create task
	models.DB.Create(&models.Task{
		Status:        models.TASK_CREATED,
		FuzzerID:      req.FuzzerID,
		CorpusID:      req.CorpusID,
		TargetID:      req.TargetID,
		MaxTime:       req.MaxTime,
		FuzzCycleTime: req.FuzzCycleTime,
	})

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

func (tc *TaskController) StopFuzz(c *gin.Context) {
	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exists")
		return
	}
	if task.Status == models.TASK_RUNNING {
		service.StopFuzz()
		models.DB.Model(task).Update("Status", models.TASK_STOPPED)
		c.JSON(http.StatusNoContent, "")
	} else {
		utils.BadRequest(c)
		return
	}
}

func (tc *TaskController) StartFuzz(c *gin.Context) {
	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exists")
		return
	}
	if task.Status == models.TASK_CREATED {
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
		ids := []uint64{task.FuzzerID, task.TargetID, task.CorpusID}
		paths := []string{}
		for i, _ := range ids {
			storageItem, err := models.GetStorageItemByID(ids[i])
			if err != nil {
				utils.DBError(c)
				return
			}
			if !utils.IsPathExists(storageItem.Path) {
				utils.InternalErrorWithMsg(c, "storageItem file missing, this should not happen")
				return
			}
			paths = append(paths, storageItem.Path)
		}

		service.Fuzz(paths[0], paths[1], paths[2], task.MaxTime, int(task.FuzzCycleTime), arguments, environments)

		models.DB.Model(task).Update("Status", models.TASK_RUNNING)
		c.JSON(http.StatusNoContent, "")
	} else {
		utils.BadRequest(c)
		return
	}
}

func (tc *TaskController) Destroy(c *gin.Context) {
	service.StopFuzz()
	_, err := models.GetTask()
	if err != nil {
		utils.DBError(c)
		return
	}
	models.DB.Delete(&models.Task{})
	models.DB.Delete(&models.TaskArgument{})
	models.DB.Delete(&models.TaskEnvironment{})
	c.JSON(http.StatusNoContent, "")
}
