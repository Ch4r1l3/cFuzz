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

// swagger:model
type TaskCreateReq struct {
	// example: test
	// required: true
	Name string `json:"name" binding:"required"`

	// example: afl-image
	Image string `json:"image"`

	// example: 1
	DeploymentID uint64 `json:"deploymentid"`

	// example: 3600
	// required: true
	Time uint64 `json:"time" binding:"required"`

	// example: 60
	// required: true
	FuzzCycleTime uint64 `json:"fuzzCycleTime" binding:"required"`

	// example: 1
	FuzzerID uint64 `json:"fuzzerID"`

	// example: 2
	CorpusID uint64 `json:"corpusID"`

	// example: 3
	TargetID uint64 `json:"targetID"`

	// example: ["AFL_FUZZ=1", "ASAN=1"]
	Environments []string `json:"environments"`

	// example: {"MEMORY_LIMIT": "100"}
	Arguments map[string]string `json:"arguments"`
}

// swagger:model
type TaskUpdateReq struct {
	// example: test
	Name string `json:"name"`

	// example: afl-image
	Image string `json:"image"`

	// example: 1
	DeploymentID uint64 `json:"deploymentid"`

	// example: 3600
	Time uint64 `json:"time"`

	// example: 60
	FuzzCycleTime uint64 `json:"fuzzCycleTime"`

	// example: 1
	FuzzerID uint64 `json:"fuzzerID"`

	// example: 2
	CorpusID uint64 `json:"corpusID"`

	// example: 3
	TargetID uint64 `json:"targetID"`

	// example: ["AFL_FUZZ=1", "ASAN=1"]
	Environments []string `json:"environments"`

	// example: {"MEMORY_LIMIT": "100"}
	Arguments map[string]string `json:"arguments"`
}

// swagger:model
type TaskResp struct {
	models.Task

	// example: ["AFL_FUZZ=1", "ASAN=1"]
	Environments []string `json:"environments"`

	// example: {"MEMORY_LIMIT": "100"}
	Arguments map[string]string `json:"arguments"`
}

type TaskController struct{}

// list tasks
func (tc *TaskController) List(c *gin.Context) {
	// swagger:operation GET /task task listTask
	// list tasks
	//
	// list tasks
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: offset
	//   in: query
	//   type: integer
	// - name: limit
	//   in: query
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        type: array
	//        items:
	//          "$ref": "#/definitions/TaskResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var tasks []models.Task
	var err error
	if !c.GetBool("pagination") {
		err = models.GetObjects(&tasks)
	} else {
		offset := c.GetInt("offset")
		limit := c.GetInt("limit")
		err = models.GetObjectsPagination(&tasks, offset, limit)
	}
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
		results = append(results, TaskResp{
			Task:         task,
			Environments: environments,
			Arguments:    arguments,
		})
	}
	c.JSON(http.StatusOK, results)
}

// retrieve task
func (tc *TaskController) Retrieve(c *gin.Context, id uint64) {
	// swagger:operation GET /task/{id} task retrieveTask
	// retrieve task
	//
	// retrieve task
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/TaskResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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
	c.JSON(http.StatusOK, TaskResp{
		Task:         task,
		Environments: environments,
		Arguments:    arguments,
	})
}

// create task
func (tc *TaskController) Create(c *gin.Context) {
	// swagger:operation POST /task task createTask
	// create task
	//
	// create task
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: taskCreateReq
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/TaskCreateReq"
	//
	// responses:
	//   '201':
	//      schema:
	//        "$ref": "#/definitions/TaskResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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
	c.JSON(http.StatusCreated, task)
}

// start task
func (tc *TaskController) Start(c *gin.Context) {
	// swagger:operation POST /task/{id}/start task startTask
	// start task
	//
	// start task
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '202':
	//     description: start task success
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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
	c.JSON(http.StatusAccepted, "")
}

// stop task
func (tc *TaskController) Stop(c *gin.Context) {
	// swagger:operation POST /task/{id}/stop task stopTask
	// stop task
	//
	// stop task
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '202':
	//     description: stop task success
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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
	c.JSON(http.StatusAccepted, "")
	return
}

// update task
func (tc *TaskController) Update(c *gin.Context) {
	// swagger:operation PUT /task/{id} task updateTask
	// update task
	//
	// update task
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: taskUpdateReq
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/TaskUpdateReq"
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '201':
	//      description: update task success
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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

	c.JSON(http.StatusCreated, "")
}

// delete task
func (tc *TaskController) Destroy(c *gin.Context) {
	// swagger:operation DELETE /task/{id} task deleteTask
	// delete task
	//
	// delete task
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '204':
	//      description: update task success
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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
