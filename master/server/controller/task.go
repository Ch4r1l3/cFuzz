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
)

type TaskCreateReq struct {
	Name          string            `json:"name" binding:"required"`
	Image         string            `json:"image"`
	DeploymentID  uint64            `json:"deploymentid"`
	Time          uint64            `json:"time" binding:"required"`
	FuzzCycleTime uint64            `json:"fuzzCycleTime" binding:"required"`
	FuzzerID      uint64            `json:"fuzzerid" binding:"required"`
	Environments  []string          `json:"environments"`
	Arguments     map[string]string `json:"arguments"`
}

type TaskUpdateReq struct {
	Name          string            `json:"name"`
	Image         string            `json:"image"`
	DeploymentID  uint64            `json:"deploymentid"`
	Time          uint64            `json:"time"`
	FuzzCycleTime uint64            `json:"fuzzCycleTime"`
	FuzzerID      uint64            `json:"fuzzerid"`
	Environments  []string          `json:"environments"`
	Arguments     map[string]string `json:"arguments"`
	Status        string            `json:"status"`
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
			"id":            task.ID,
			"deploymentid":  task.DeploymentID,
			"name":          task.Name,
			"image":         task.Image,
			"time":          task.Time,
			"fuzzCycleTime": task.FuzzCycleTime,
			"fuzzerid":      task.FuzzerID,
			"status":        task.Status,
			"errorMsg":      task.ErrorMsg,
			"environments":  environments,
			"arguments":     arguments,
		})
	}
	c.JSON(http.StatusOK, results)
}

func (tc *TaskController) Retrieve(c *gin.Context, id uint64) {
	var task models.Task
	err := models.GetObjectByID(&task, id)
	if err != nil {
		utils.DBError(c)
		return
	}
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
	result := map[string]interface{}{
		"id":            task.ID,
		"deploymentid":  task.DeploymentID,
		"name":          task.Name,
		"image":         task.Image,
		"time":          task.Time,
		"fuzzCycleTime": task.FuzzCycleTime,
		"fuzzerid":      task.FuzzerID,
		"status":        task.Status,
		"errorMsg":      task.ErrorMsg,
		"environments":  environments,
		"arguments":     arguments,
	}
	c.JSON(http.StatusOK, result)
}

func (tc *TaskController) Create(c *gin.Context) {
	var req TaskCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	if req.Image == "" && req.DeploymentID == 0 {
		utils.BadRequestWithMsg(c, "image and deployment is empty")
		return
	}
	task := models.Task{
		FuzzerID:      req.FuzzerID,
		FuzzCycleTime: req.FuzzCycleTime,
		Time:          req.Time,
		Name:          req.Name,
		Status:        models.TaskCreated,
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

func (tc *TaskController) taskStart(c *gin.Context, task *models.Task) {
	var Err error
	if !models.IsObjectExistsByTaskID(&models.TaskTarget{}, task.ID) {
		utils.BadRequestWithMsg(c, "you should upload target first")
		return
	}
	if !models.IsObjectExistsByTaskID(&models.TaskCorpus{}, task.ID) {
		utils.BadRequestWithMsg(c, "you should upload corpus first")
		return
	}

	err := service.CreateServiceByTaskID(task.ID)
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
		Err = err
		utils.InternalErrorWithMsg(c, "create deployment failed")
		return
	}
	if err = models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Status", models.TaskStarted).Error; err != nil {
		service.DeleteDeployByTaskID(task.ID)
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}

func (tc *TaskController) taskStop(c *gin.Context, taskID uint64) {
	if err := models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskStopped).Error; err != nil {
		utils.DBError(c)
		return
	}
	err1 := service.DeleteServiceByTaskID(taskID)
	err2 := service.DeleteDeployByTaskID(taskID)
	if err1 != nil || err2 != nil {
		utils.InternalErrorWithMsg(c, "kubernetes delete error")
		return
	}
	c.JSON(http.StatusOK, "")
	return
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
	if task.Status == models.TaskCreated && req.Status == models.TaskStarted {
		tc.taskStart(c, &task)
		return
	} else if task.Status != models.TaskCreated && task.Status != models.TaskStopped && task.Status != models.TaskError && req.Status == models.TaskStopped {
		tc.taskStop(c, task.ID)
		return
	} else if task.Status != models.TaskCreated {
		utils.BadRequestWithMsg(c, "task already started or stopped")
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
	if req.FuzzCycleTime != 0 {
		if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update("FuzzCycleTime", req.Time).Error; err != nil {
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
