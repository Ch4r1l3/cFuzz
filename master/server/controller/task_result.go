package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskResultController struct{}

func (trc *TaskResultController) Retrieve(c *gin.Context, taskID uint64) {
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
	c.JSON(http.StatusOK, gin.H{
		"command":      taskFuzzResult[0].Command,
		"timeExecuted": taskFuzzResult[0].TimeExecuted,
		"updateAt":     taskFuzzResult[0].UpdateAt,
		"stats":        stats,
	})

}
