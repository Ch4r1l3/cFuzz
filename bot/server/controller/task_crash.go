package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskCrashController struct{}

// List Crashes
func (tcc *TaskCrashController) List(c *gin.Context) {
	// swagger:operation GET /task/crash taskCrash listTaskCrash
	// list all crash
	//
	// list all crash
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//      schema:
	//        type: array
	//        items:
	//          "$ref": "#/definitions/TaskCrash"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	crashes, err := models.GetCrashes()
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, crashes)
}

// Download Crash
func (tcc *TaskCrashController) Download(c *gin.Context) {
	// swagger:operation GET /task/crash/{id} taskCrash downloadTaskCrash
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

	var req TaskIDReq
	if err := c.ShouldBindUri(&req); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
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
