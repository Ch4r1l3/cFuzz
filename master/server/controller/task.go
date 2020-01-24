package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"net/http"
	"time"
)

type TaskCreateReq struct {
	Name          string            `json:"name" binding:"required"`
	Image         string            `json:"image"`
	DeploymentID  uint64            `json:"deploymentid"`
	Time          uint64            `json:"time" binding:"required"`
	FuzzCycleTime uint64            `json:"fuzzCycleTime" binding:"required"`
	FuzzerID      uint64            `json:"fuzzerid"`
	CorpusID      uint64            `json:"corpusid"`
	TargetID      uint64            `json:"targetid"`
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
	CorpusID      uint64            `json:"corpusid"`
	TargetID      uint64            `json:"targetid"`
	Environments  []string          `json:"environments"`
	Arguments     map[string]string `json:"arguments"`
	Status        string            `json:"status"`
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
			"corpusid":      task.CorpusID,
			"targetid":      task.TargetID,
			"status":        task.Status,
			"errorMsg":      task.ErrorMsg,
			"startedAt":     task.StartedAt,
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
		"corpusid":      task.CorpusID,
		"targetid":      task.TargetID,
		"status":        task.Status,
		"errorMsg":      task.ErrorMsg,
		"startedAt":     task.StartedAt,
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
		CorpusID:      req.CorpusID,
		TargetID:      req.TargetID,
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
	if req.FuzzerID != 0 && !models.IsObjectExistsByID(&models.StorageItem{}, req.FuzzerID) {
		utils.BadRequestWithMsg(c, "fuzzer not exists")
		return
	}
	if req.CorpusID != 0 && !models.IsObjectExistsByID(&models.StorageItem{}, req.CorpusID) {
		utils.BadRequestWithMsg(c, "corpus not exists")
		return
	}
	if req.TargetID != 0 && !models.IsObjectExistsByID(&models.StorageItem{}, req.TargetID) {
		utils.BadRequestWithMsg(c, "target not exists")
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

func (tc *TaskController) Start(c *gin.Context) {
	var uriReq UriIDReq
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
	if task.Status != models.TaskCreated {
		utils.BadRequestWithMsg(c, "wrong status")
		return
	}
	var Err error
	if !models.IsObjectExistsByID(&models.StorageItem{}, task.FuzzerID) {
		utils.BadRequestWithMsg(c, "you should upload fuzzer first")
		return
	}
	if !models.IsObjectExistsByID(&models.StorageItem{}, task.TargetID) {
		utils.BadRequestWithMsg(c, "you should upload target first")
		return
	}
	if !models.IsObjectExistsByID(&models.StorageItem{}, task.CorpusID) {
		utils.BadRequestWithMsg(c, "you should upload corpus first")
		return
	}

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
		Err = err
		utils.InternalErrorWithMsg(c, "create deployment failed")
		return
	}
	if err = models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Status", models.TaskStarted).Error; err != nil {
		service.DeleteDeployByTaskID(task.ID)
		utils.DBError(c)
		Err = err
		return
	}
	if err = models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("StatusUpdateAt", time.Now().Unix()).Error; err != nil {
		service.DeleteDeployByTaskID(task.ID)
		utils.DBError(c)
		Err = err
		return
	}
	c.JSON(http.StatusNoContent, "")
}

func (tc *TaskController) Stop(c *gin.Context) {
	var uriReq UriIDReq
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
	if task.Status != models.TaskStarted && task.Status != models.TaskInitializing && task.Status != models.TaskRunning {
		utils.BadRequestWithMsg(c, "wrong status")
		return
	}
	if err := models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Status", models.TaskStopped).Error; err != nil {
		utils.DBError(c)
		return
	}
	err1 := service.DeleteServiceByTaskID(task.ID)
	err2 := service.DeleteDeployByTaskID(task.ID)
	if err1 != nil || err2 != nil {
		utils.InternalErrorWithMsg(c, "kubernetes delete error")
		return
	}
	c.JSON(http.StatusNoContent, "")
	return
}

func (tc *TaskController) Update(c *gin.Context) {
	var uriReq UriIDReq
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

	if task.Status != models.TaskCreated {
		utils.BadRequestWithMsg(c, "you can change task after it started")
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
	ids := []uint64{req.FuzzerID, req.CorpusID, req.TargetID}
	types := []string{models.Fuzzer, models.Corpus, models.Target}
	modelsField := []string{"FuzzerID", "CorpusID", "TargetID"}
	for i, _ := range ids {
		if ids[i] != 0 {
			var storageItem models.StorageItem
			err = models.GetObjectByID(&storageItem, ids[i])
			if err != nil {
				utils.BadRequest(c)
				return
			}
			if storageItem.Type != types[i] {
				utils.BadRequestWithMsg(c, "wrong type")
				return
			}
			if err = models.DB.Model(&models.Task{}).Where("id = ?", uriReq.ID).Update(modelsField[i], ids[i]).Error; err != nil {
				utils.DBError(c)
				return
			}
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

	c.JSON(http.StatusNoContent, "")
}

func (tc *TaskController) Destroy(c *gin.Context) {
	var uriReq UriIDReq
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
	if err := models.DeleteTask(uriReq.ID); err != nil {
		utils.DBError(c)
		return
	}
	service.DeleteServiceByTaskID(uriReq.ID)
	service.DeleteDeployByTaskID(uriReq.ID)
	c.JSON(http.StatusNoContent, "")
}
