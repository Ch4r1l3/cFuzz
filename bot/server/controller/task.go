package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/bot/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskController struct{}

// swagger:model
type TaskCreateReq struct {
	//example: 1
	//required: true
	FuzzerID uint64 `json:"fuzzerID" binding:"required"`

	//example: 2
	//required: true
	CorpusID uint64 `json:"corpusID" binding:"required"`

	//example: 3
	//required: true
	TargetID uint64 `json:"targetID" binding:"required"`

	//example: 3600
	//required: true
	MaxTime int `json:"maxTime" binding:"required"`

	//example: 60
	//required: true
	FuzzCycleTime uint64 `json:"fuzzCycleTime" binding:"required"`

	Arguments map[string]string `json:"arguments"`

	//example: [ASAN_ON=true, ASAN_AFL=true]
	Environments []string `json:"environments"`
}

type TaskIDReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

// Retrieve Task
func (tc *TaskController) Retrieve(c *gin.Context) {
	// swagger:operation GET /task task retrieveTask
	// retrieve task
	//
	// retrieve task
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/TaskCreateReq"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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

// Create Task
func (tc *TaskController) Create(c *gin.Context) {
	// swagger:operation POST /task task createTask
	// retrieve task
	//
	// retrieve task
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: TaskCreateReq
	//   in: body
	//   required: true
	//   schema:
	//       "$ref": "#/definitions/TaskCreateReq"
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/TaskCreateReq"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	//var task models.Task
	var req TaskCreateReq
	err := c.BindJSON(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	task, err := models.GetTask()
	if err == nil && task.Status == models.TaskRunning {
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
		Status:        models.TaskCreated,
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

// Stop Fuzz
func (tc *TaskController) StopFuzz(c *gin.Context) {
	// swagger:operation POST /task/stop task stopTask
	// stop task
	//
	// stop task
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
	//      description: "stop task success"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exists")
		return
	}
	if task.Status == models.TaskRunning {
		service.StopFuzz()
		models.DB.Model(task).Update("Status", models.TaskStopped)
		c.JSON(http.StatusNoContent, "")
	} else {
		utils.BadRequestWithMsg(c, "task is not running")
	}
}

// Start Fuzz
func (tc *TaskController) StartFuzz(c *gin.Context) {
	// swagger:operation POST /task/start task startTask
	// start task
	//
	// start task
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
	//      description: "start task success"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	task, err := models.GetTask()
	if err != nil {
		utils.BadRequestWithMsg(c, "task not exists")
		return
	}
	if task.Status == models.TaskCreated {
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

		models.DB.Model(task).Update("Status", models.TaskRunning)
		c.JSON(http.StatusNoContent, "")
	} else {
		utils.BadRequestWithMsg(c, "task status is not created")
	}
}

// Delete Task
func (tc *TaskController) Destroy(c *gin.Context) {
	// swagger:operation DELETE /task task deleteTask
	// delete task
	//
	// delete task
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
	//      description: "delete task success"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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
