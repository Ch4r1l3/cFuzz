package controller

import (
	"errors"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
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
	ID uint64 `json:"id"`
}

type TaskIDUriReq struct {
	TaskID uint64 `json:"taskid"`
}

func getTaskID(c *gin.Context) (uint64, error) {
	var req TaskIDUriReq
	if err := c.ShouldBindUri(&req); err != nil {
		utils.BadRequest(c)
		return 0, err
	}
	if !models.IsObjectExistsByID(&models.Task{}, req.TaskID) {
		utils.NotFound(c)
		return 0, errors.New("not exists")
	}
	return req.TaskID, nil
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
	} else if p1 != "" && p2 != "" && p3 == "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		if p2 == "target" {
			taskTarget := new(TaskTargetController)
			taskTarget.Destroy(c, n)
		} else if p2 == "corpus" {
			taskCorpus := new(TaskCorpusController)
			taskCorpus.Destroy(c, n)
		}
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
		taskCorpus.DestroyByID(c, n1, n2)
	} else {
		utils.BadRequest(c)
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
	results := []interface{}{}
	for _, task := range tasks {
		environments, err := models.GetEnvironments(task.ID)
		if err != nil {
			utils.DBError(c)
			return
		}
		arguments, err := models.GetArguments(task.ID)
		if err != nil {
			utils.DBError(c)
			return
		}
		results = append(results, map[string]interface{}{
			"id":           task.ID,
			"dockerfileid": task.DockerfileID,
			"time":         task.Time,
			"fuzzerid":     task.FuzzerID,
			"running":      task.Running,
			"environments": environments,
			"arguments":    arguments,
		})
	}
	c.JSON(http.StatusOK, results)
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
			models.DeleteTask(task.ID)
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
		fmt.Println(req.Arguments)
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
	var task models.Task
	if err = models.GetObjectByID(&task, uriReq.ID); err != nil {
		utils.NotFound(c)
		return
	}
	var req TaskUpdateReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	if task.Running && req.Running {
		utils.BadRequest(c)
		return
	}
	if !task.Running && req.Running {
		if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("Running", req.Running).Error; err != nil {
			utils.DBError(c)
			return
		}
		c.JSON(http.StatusOK, "")
		return
	}
	if task.Running && !req.Running {
		if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("Running", req.Running).Error; err != nil {
			utils.DBError(c)
			return
		}
	}
	if req.DockerfileID != 0 {
		if models.IsObjectExistsByID(&models.Dockerfile{}, req.DockerfileID) {
			if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("DockerfileID", req.DockerfileID).Error; err != nil {
				utils.DBError(c)
				return
			}
		} else {
			utils.BadRequest(c)
			return
		}
	}
	if req.FuzzerID != 0 {
		if models.IsObjectExistsByID(&models.Fuzzer{}, req.FuzzerID) {
			if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("FuzzerID", req.FuzzerID).Error; err != nil {
				utils.DBError(c)
				return
			}
		} else {
			utils.BadRequest(c)
			return
		}
	}
	if req.Time != 0 {
		if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("Time", req.Time).Error; err != nil {
			utils.DBError(c)
			return
		}
	}
	if req.Arguments != nil {
		if err = models.DeleteObjectsByTaskID(&models.TaskArgument{}, uriReq.ID); err != nil {
			utils.DBError(c)
			return
		}
		if err = models.InsertArguments(uriReq.ID, req.Arguments); err != nil {
			utils.DBError(c)
			return
		}
	}
	if req.Environments != nil {
		if err = models.DeleteObjectsByTaskID(&models.TaskEnvironment{}, uriReq.ID); err != nil {
			utils.DBError(c)
			return
		}
		if err = models.InsertEnvironments(uriReq.ID, req.Environments); err != nil {
			utils.DBError(c)
			return
		}
	}

	c.JSON(http.StatusOK, "")
}

func (tc *TaskController) Destroy(c *gin.Context, id uint64) {
	if err := models.DeleteTask(id); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}

type TaskCorpusController struct{}

func (tcc *TaskCorpusController) List(c *gin.Context) {
	var taskid uint64
	var err error
	if taskid, err = getTaskID(c); err != nil {
		return
	}
	var corpus []models.TaskCorpus
	if err := models.GetObjectsByTaskID(corpus, taskid); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, corpus)
}

func (tcc *TaskCorpusController) Create(c *gin.Context) {
	var taskid uint64
	var err error
	if taskid, err = getTaskID(c); err != nil {
		return
	}
	var tempFile string
	if tempFile, err = utils.SaveTempFile(c, "file", "corpus"); err != nil {
		return
	}
	corpus := models.TaskCorpus{
		TaskID:   taskid,
		Path:     tempFile,
		FileName: filepath.Base(tempFile),
	}
	if err := models.InsertObject(&corpus); err != nil {
		os.RemoveAll(tempFile)
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, corpus)
}

func (tcc *TaskCorpusController) Destroy(c *gin.Context, taskid uint64) {
	var corpus []models.TaskCorpus
	if err := models.GetObjectsByTaskID(corpus, taskid); err != nil {
		utils.DBError(c)
		return
	}
	for _, v := range corpus {
		os.RemoveAll(v.Path)
	}
	if err := models.DeleteObjectsByTaskID(&models.TaskCorpus{}, taskid); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}

func (tcc *TaskCorpusController) DestroyByID(c *gin.Context, taskid uint64, corpusid uint64) {
	if !models.IsObjectExistsByID(&models.Task{}, taskid) {
		utils.NotFound(c)
		return
	}
	if !models.IsObjectExistsByID(&models.TaskCorpus{}, corpusid) {
		utils.NotFound(c)
		return
	}
	if err := models.DeleteObjectByID(&models.TaskCorpus{}, corpusid); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}

type TaskTargetController struct{}

func (ttc *TaskTargetController) Retrieve(c *gin.Context) {
	var taskid uint64
	var err error
	if taskid, err = getTaskID(c); err != nil {
		return
	}
	var targets []models.TaskTarget
	if err = models.GetObjectsByTaskID(targets, taskid); err != nil {
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
	var tempFile string
	if tempFile, err = utils.SaveTempFile(c, "file", "target"); err != nil {
		return
	}
	corpus := models.TaskTarget{
		TaskID: taskid,
		Path:   tempFile,
	}
	if err = models.InsertObject(&corpus); err != nil {
		os.RemoveAll(tempFile)
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, corpus)

}

func (ttc *TaskTargetController) Destroy(c *gin.Context, id uint64) {
	var taskid uint64
	var err error
	if taskid, err = getTaskID(c); err != nil {
		return
	}
	var targets []models.TaskTarget
	if err = models.GetObjectsByTaskID(targets, taskid); err != nil {
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
