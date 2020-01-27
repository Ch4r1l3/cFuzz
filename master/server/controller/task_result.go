package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
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

	var taskFuzzResult []models.TaskFuzzResult
	if err := models.GetObjectsByTaskID(&taskFuzzResult, taskID); err != nil {
		utils.DBError(c)
		return
	}
	if len(taskFuzzResult) == 0 {
		c.JSON(http.StatusOK, "")
		return
	}
	stats, err := models.GetTaskFuzzResultStat(taskFuzzResult[0].ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, TaskResultResp{
		TaskFuzzResult: taskFuzzResult[0],
		Stats:          stats,
	})
}
