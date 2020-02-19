package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// swagger:model
type TaskResultResp struct {
	models.TaskFuzzResult
	// example: {"crashes": "1"}
	// in: body
	Stats map[string]string `json:"stats"`
}

type TaskResultController struct{}

// Retrieve Task Result
func (trc *TaskResultController) Retrieve(c *gin.Context, taskID uint64) {
	// swagger:operation GET /task/{taskID}/result taskResult retrieveTaskResult
	// retrieve task result
	//
	// retrieve task result
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: taskID
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/TaskResultResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var task models.Task
	if err := service.GetObjectByID(&task, taskID); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.NotFound(c)
			return
		}
		utils.DBError(c)
		return
	}
	if task.UserID != uint64(c.GetInt64("id")) && !c.GetBool("isAdmin") {
		utils.Forbidden(c)
		return
	}
	taskFuzzResult, err := service.GetLastestFuzzResultByTaskID(taskID)
	if err != nil {
		utils.DBError(c)
		return
	}
	if taskFuzzResult == nil {
		c.String(http.StatusOK, "")
		return
	}
	stats, err := service.GetTaskFuzzResultStat(taskFuzzResult.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, TaskResultResp{
		TaskFuzzResult: *taskFuzzResult,
		Stats:          stats,
	})
}
