package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskCrashController struct{}

func (tcc *TaskCrashController) List(c *gin.Context, taskID uint64) {
	var crashes []models.TaskCrash
	err := models.GetObjectsByTaskID(&crashes, taskID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, crashes)
}

func (tcc *TaskCrashController) Download(c *gin.Context, taskID uint64, crashID uint64) {
	var crash models.TaskCrash
	err := models.GetObjectByTaskIDAndID(&crash, taskID, crashID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=crash%d", crashID))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(crash.Path)
}
