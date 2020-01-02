package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskCreateReq struct {
	DockerfileID uint64            `json:"dockerfileid" binding:"required"`
	Time         uint64            `json:"time" binding:"required"`
	FuzzerID     uint64            `json:"fuzzerid" binding:"required"`
	Environments []string          `json:"environments"`
	Arguments    map[string]string `json:"arguments"`
}

type TaskUpdateReq struct {
	DockerfileID uint64            `json:"dockerfileid"`
	Time         uint64            `json:"time"`
	FuzzerID     uint64            `json:"fuzzerid"`
	Environments []string          `json:"environments"`
	Arguments    map[string]string `json:"arguments"`
	Running      bool              `json:"running"`
}

type TaskUpdateUriReq struct {
	id uint64 `json:"id"`
}

func TaskDeleteHandler(c *gin.Context) {
	p1 := c.Param("path1")
	p2 := c.Param("path2")
	p3 := c.Param("path3")
	if p1 != "" && p2 == "" && p3 == "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		task := new(TaskController)
		task.Destroy(c, n)
	} else if p1 != "" && p2 == "target" && p3 == "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}

		taskTarget := new(TaskTargetController)
		taskTarget.Destroy(c, n)
	} else if p1 != "" && p2 == "corpus" && p3 != "" {
		n1, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}

		n2, err := strconv.ParseUint(p3, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		taskCorpus := new(TaskCorpusController)
		taskCorpus.Destroy(c, n1, n2)
	}
}

type TaskController struct{}

func (tc *TaskController) List(c *gin.Context) {
	tasks := []models.Task{}
	err := models.GetObjects(&tasks)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) Create(c *gin.Context) {
	var req TaskCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	if !models.IsDockerfileExistsByID(req.DockerfileID) || !models.IsFuzzerExistsByID(req.FuzzerID) {
		utils.BadRequestWithMsg(c, "dockerfile not exists or fuzzer not exists")
		return
	}
	task := models.Task{
		DockerfileID: req.DockerfileID,
		FuzzerID:     req.FuzzerID,
		Time:         req.Time,
	}
	err = models.InsertObject(&task)
	if err != nil {
		utils.DBError(c)
		return
	}
	var Err error
	defer func() {
		if Err != nil {
			models.DeleteObjectsByTaskID(&models.TaskEnvironment{}, task.ID)
			models.DeleteObjectsByTaskID(&models.TaskArgument{}, task.ID)
			models.DeleteObjectByID(&models.Task{}, task.ID)
		}
	}()
	if req.Environments != nil {
		fmt.Println(req.Environments)
		err = models.InsertEnvironments(task.ID, req.Environments)
		if err != nil {
			Err = err
			utils.DBError(c)
			return
		}
	}
	if req.Arguments != nil {
		err = models.InsertArguments(task.ID, req.Arguments)
		if err != nil {
			Err = err
			utils.DBError(c)
			return
		}

	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) Update(c *gin.Context) {
	var uriReq TaskUpdateUriReq
	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	var req TaskUpdateReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	if req.DockerfileID != 0 {
		if models.IsObjectExistsByID(&models.Dockerfile{}, req.DockerfileID) {

		} else {
			utils.BadRequest(c)
			return
		}
	}
}

func (tc *TaskController) Destroy(c *gin.Context, id uint64) {
	if err := models.DeleteObjectByID(&models.Task{}, id); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")

}

type TaskCorpusController struct{}

func (tcc *TaskCorpusController) List(c *gin.Context) {

}

func (tcc *TaskCorpusController) Create(c *gin.Context) {

}

func (tcc *TaskCorpusController) Destroy(c *gin.Context, taskid uint64, corpusid uint64) {

}

type TaskTargetController struct{}

func (ttc *TaskTargetController) Retrieve(c *gin.Context) {

}

func (ttc *TaskTargetController) Create(c *gin.Context) {

}

func (ttc *TaskTargetController) Destroy(c *gin.Context, id uint64) {

}
