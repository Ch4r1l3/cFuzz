package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskCrashController struct{}

// swagger:model
type TaskCrashListResp struct {
	Data []models.TaskCrash `json:"data"`
	CountResp
}

// List all crashes by taskID
func (tcc *TaskCrashController) ListByTaskID(c *gin.Context, taskID uint64) {
	// swagger:operation GET /task/{taskID}/crash taskCrash listTaskCrash
	// list all crashes
	//
	// list all crashes
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: taskID
	//   in: path
	//   required: true
	//   type: integer
	// - name: offset
	//   in: query
	//   type: integer
	// - name: limit
	//   in: query
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/TaskCrashListResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var crashes []models.TaskCrash
	task, err := service.GetTaskByID(taskID)
	if err != nil {
		utils.InternalErrorWithMsg(c, err.Error())
		return
	}
	if task == nil {
		utils.NotFound(c)
		return
	}
	if task.UserID != uint64(c.GetInt64("id")) && !c.GetBool("isAdmin") {
		utils.Forbidden(c)
		return
	}
	offset := c.GetInt("offset")
	limit := c.GetInt("limit")
	count, err := service.GetObjectsByTaskIDPagination(&crashes, taskID, offset, limit)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, TaskCrashListResp{
		Data: crashes,
		CountResp: CountResp{
			Count: count,
		},
	})
}

// Download task crash
func (tcc *TaskCrashController) Download(c *gin.Context) {
	// swagger:operation GET /crash/{id} taskCrash downloadCrash
	// download crash by id
	//
	// download crash by id
	// ---
	// produces:
	// - application/octet-stream
	//
	// parameters:
	// - name: id
	//   description: id of crash
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        type: file
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var crash models.TaskCrash
	err := getObject(c, &crash)
	if err != nil {
		return
	}
	task, err := service.GetTaskByID(crash.TaskID)
	if err != nil {
		utils.DBError(c)
		return
	}
	if task == nil {
		utils.InternalErrorWithMsg(c, "task not exists")
	}
	if task.UserID != uint64(c.GetInt64("id")) && !c.GetBool("isAdmin") {
		utils.Forbidden(c)
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=crash%d", crash.ID))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(crash.Path)
}
