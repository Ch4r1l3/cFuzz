package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type TaskTargetController struct{}

func (ttc *TaskTargetController) Retrieve(c *gin.Context, taskid uint64) {
	var targets []models.TaskTarget
	if err := models.GetObjectsByTaskID(&targets, taskid); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, targets)
}

func (ttc *TaskTargetController) Create(c *gin.Context) {
	var taskid uint64
	var err error
	if taskid, err = getTaskID(c); err != nil {
		return
	}
	var targets []models.TaskTarget
	if err = models.GetObjectsByTaskID(&targets, taskid); err != nil {
		utils.DBError(c)
		return
	}
	if len(targets) > 0 {
		utils.BadRequestWithMsg(c, "you should delete target first")
		return
	}
	var tempFile string
	if tempFile, err = utils.SaveTempFile(c, "file", "target"); err != nil {
		return
	}
	target := models.TaskTarget{
		TaskID: taskid,
		Path:   tempFile,
	}
	if err = models.InsertObject(&target); err != nil {
		os.RemoveAll(tempFile)
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, target)
}

func (ttc *TaskTargetController) Destroy(c *gin.Context, id uint64) {
	var taskid uint64
	var err error
	if taskid, err = getTaskID(c); err != nil {
		return
	}
	var targets []models.TaskTarget
	if err = models.GetObjectsByTaskID(&targets, taskid); err != nil {
		utils.DBError(c)
		return
	}
	for _, v := range targets {
		os.RemoveAll(v.Path)
	}
	if err = models.DeleteObjectsByTaskID(&models.TaskTarget{}, taskid); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")

}
