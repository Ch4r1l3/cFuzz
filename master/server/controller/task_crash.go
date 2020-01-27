package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskCrashController struct{}

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

	var crashes []models.TaskCrash
	err := models.GetObjectsByTaskID(&crashes, taskID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, crashes)
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

	var req UriIDReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	var crash models.TaskCrash
	err = models.GetObjectByID(&crash, req.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=crash%d", req.ID))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(crash.Path)
}
