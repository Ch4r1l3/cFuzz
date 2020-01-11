package controller

import (
	"errors"
	//"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type TaskCreateReq struct {
	Name         string            `json:"name" binding:"required"`
	Image        string            `json:"image"`
	DeploymentID uint64            `json:"deploymentid"`
	Time         uint64            `json:"time" binding:"required"`
	FuzzerID     uint64            `json:"fuzzerid" binding:"required"`
	Environments []string          `json:"environments"`
	Arguments    map[string]string `json:"arguments"`
}

type TaskUpdateReq struct {
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	DeploymentID uint64            `json:"deploymentid"`
	Time         uint64            `json:"time"`
	FuzzerID     uint64            `json:"fuzzerid"`
	Environments []string          `json:"environments"`
	Arguments    map[string]string `json:"arguments"`
	Running      bool              `json:"running"`
}

type TaskUpdateUriReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

type TaskIDUriReq struct {
	TaskID uint64 `uri:"taskid" binding:"required"`
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
			"deploymentid": task.DeploymentID,
			"name":         task.Name,
			"image":        task.Image,
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
	if req.Image == "" && req.DeploymentID == 0 {
		utils.BadRequest(c)
		return
	}
	task := models.Task{
		FuzzerID: req.FuzzerID,
		Time:     req.Time,
		Name:     req.Name,
	}
	if req.Image != "" {
		task.Image = req.Image
	} else {
		if !models.IsObjectExistsByID(&models.Deployment{}, req.DeploymentID) {
			utils.BadRequestWithMsg(c, "deployment not exists")
		}
		task.DeploymentID = req.DeploymentID
	}
	if !models.IsFuzzerExistsByID(req.FuzzerID) {
		utils.BadRequestWithMsg(c, "fuzzer not exists")
		return
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
	var Err error
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
		err = service.CreateServiceByTaskID(task.ID)
		if err != nil {
			utils.InternalErrorWithMsg(c, "create service failed")
			return
		}
		defer func() {
			if Err != nil {
				service.DeleteServiceByTaskID(task.ID)
			}
		}()
		var deployment *appsv1.Deployment
		if task.Image != "" {
			deployment, err = service.GenerateDeployment(task.ID, task.Name, task.Image, 1)
			if err != nil {
				Err = err
				utils.InternalErrorWithMsg(c, "generate deployment failed")
				return
			}
		} else if task.DeploymentID != 0 {
			var tempDeployment models.Deployment
			if err = models.GetObjectByID(&tempDeployment, task.ID); err != nil {
				Err = err
				utils.BadRequestWithMsg(c, "deployment not exists")
				return
			}
			deployment, err = service.GenerateDeploymentByYaml(tempDeployment.Content, task.ID)
			if err != nil {
				Err = err
				utils.BadRequestWithMsg(c, err.Error())
				return
			}
		} else {
			utils.BadRequestWithMsg(c, "image or deployment should have value")
			return
		}
		err = service.CreateDeploy(deployment)
		if err != nil {
			//fmt.Println(err.Error())
			Err = err
			utils.InternalErrorWithMsg(c, "create deployment failed")
			return
		}
		if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("Running", req.Running).Error; err != nil {
			service.DeleteDeployByTaskID(task.ID)
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
		err1 := service.DeleteServiceByTaskID(task.ID)
		err2 := service.DeleteDeployByTaskID(task.ID)
		if err1 != nil || err2 != nil {
			utils.InternalErrorWithMsg(c, "kubernetes delete error")
			return
		}
		c.JSON(http.StatusOK, "")
		return
	}
	if req.Image != "" {
		if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("DeploymentID", 0).Error; err != nil {
			utils.DBError(c)
			return
		}
		if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("Image", req.Image).Error; err != nil {
			utils.DBError(c)
			return
		}
	} else if req.DeploymentID != 0 {
		if models.IsObjectExistsByID(&models.Deployment{}, req.DeploymentID) {
			if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("DeploymentID", req.DeploymentID).Error; err != nil {
				utils.DBError(c)
				return
			}
			if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("Image", "").Error; err != nil {
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

func (tcc *TaskCorpusController) Retrieve(c *gin.Context) {
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
	var corpusArray []models.TaskCorpus
	if err = models.GetObjectsByTaskID(corpusArray, taskid); err != nil {
		utils.DBError(c)
		return
	}
	if len(corpusArray) > 0 {
		utils.BadRequestWithMsg(c, "you should delete corpus first")
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

/*
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
*/

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
	var targets []models.TaskTarget
	if err = models.GetObjectsByTaskID(targets, taskid); err != nil {
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
