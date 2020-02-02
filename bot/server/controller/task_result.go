package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// swagger:model
type TaskResultResp struct {
	models.TaskFuzzResult

	// example: {"cycles_done": "60"}
	Stats map[string]string `json:"stats"`
}

type TaskResultController struct{}

func (trc *TaskResultController) Retrieve(c *gin.Context) {
	// swagger:operation GET /task/result taskResult listTaskResult
	// get task result
	//
	// get task result
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/TaskResultResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	result, stats, err := models.GetFuzzResult()
	if err != nil {
		utils.DBError(c)
		return
	}
	if result == nil {
		c.JSON(http.StatusOK, "")
		return
	}
	c.JSON(http.StatusOK, TaskResultResp{
		TaskFuzzResult: *result,
		Stats:          stats,
	})
}
