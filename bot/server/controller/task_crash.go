package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskCrashController struct{}

func (tcc *TaskCrashController) List(c *gin.Context) {
	crashes, err := models.GetCrashes()
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, crashes)
}

func (tcc *TaskCrashController) Download(c *gin.Context) {
	var req TaskIDReq
	if err := c.ShouldBindUri(&req); err != nil {
		utils.BadRequest(c)
		return
	}
	crash, err := models.GetCrashByID(req.ID)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=crash%d", req.ID))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(crash.Path)
}
