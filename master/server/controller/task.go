package controller

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/Ch4r1l3/cFuzz/master/server/service/kubernetes"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"net/http"
	"time"
)

func getTask(c *gin.Context) (*models.Task, error) {
	var task models.Task
	err := getObject(c, &task)
	if err != nil {
		return nil, err
	}
	if task.UserID != uint64(c.GetInt64("id")) && !c.GetBool("isAdmin") {
		utils.Forbidden(c)
		return nil, errors.New("no permission")
	}
	return &task, nil
}

// swagger:model
type TaskCreateReq struct {
	// example: test
	// required: true
	Name string `json:"name" binding:"required"`

	// example: 1
	ImageID uint64 `json:"imageID" binding:"required"`

	// example: 3600
	// required: true
	Time uint64 `json:"time" binding:"required"`

	// example: 60
	// required: true
	FuzzCycleTime uint64 `json:"fuzzCycleTime" binding:"required"`

	// example: 1
	FuzzerID uint64 `json:"fuzzerID" binding:"required"`

	// example: 2
	CorpusID uint64 `json:"corpusID" binding:"required"`

	// example: 3
	TargetID uint64 `json:"targetID" binding:"required"`

	// example: ["AFL_FUZZ=1", "ASAN=1"]
	Environments []string `json:"environments"`

	// example: {"MEMORY_LIMIT": "100"}
	Arguments map[string]string `json:"arguments"`
}

// swagger:model
type TaskUpdateReq struct {
	// example: test
	Name string `json:"name"`

	// example: 1
	ImageID uint64 `json:"imageID"`

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

	// example: 1
	CrashNum int `json:"crashNum"`

	// example: ["AFL_FUZZ=1", "ASAN=1"]
	Environments []string `json:"environments"`

	// example: {"MEMORY_LIMIT": "100"}
	Arguments map[string]string `json:"arguments"`
}

type TaskRespCombine struct {
	Data []TaskResp `json:"data"`
	CountResp
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
	// - name: name
	//   in: query
	//   type: string
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
	count, err := getList(c, &tasks)

	if err != nil {
		utils.DBError(c)
		return
	}
	results := []TaskResp{}
	for _, task := range tasks {
		environments, err := service.GetEnvironments(task.ID)
		if err != nil {
			utils.DBError(c)
			return
		}
		arguments, err := service.GetArguments(task.ID)
		if err != nil {
			utils.DBError(c)
			return
		}
		crashNum, err := service.GetCountByTaskID(&models.TaskCrash{}, task.ID)
		if err != nil {
			utils.DBError(c)
			return
		}
		results = append(results, TaskResp{
			Task:         task,
			CrashNum:     crashNum,
			Environments: environments,
			Arguments:    arguments,
		})
	}
	c.JSON(http.StatusOK, TaskRespCombine{
		Data: results,
		CountResp: CountResp{
			Count: count,
		},
	})
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

	task, err := service.GetTaskByID(id)
	if err != nil {
		utils.DBError(c)
		return
	}
	if task == nil {
		utils.NotFound(c)
		return
	}
	environments, err := service.GetEnvironments(task.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	arguments, err := service.GetArguments(task.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	crashNum, err := service.GetCountByTaskID(&models.TaskCrash{}, task.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, TaskResp{
		Task:         *task,
		CrashNum:     crashNum,
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
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
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
	task := models.Task{
		FuzzerID:      req.FuzzerID,
		CorpusID:      req.CorpusID,
		TargetID:      req.TargetID,
		ImageID:       req.ImageID,
		FuzzCycleTime: req.FuzzCycleTime,
		Time:          req.Time,
		Name:          req.Name,
		Status:        models.TaskCreated,
		UserID:        uint64(c.GetInt64("id")),
	}
	if !service.IsObjectExistsByID(&models.Image{}, req.ImageID) {
		utils.BadRequestWithMsg(c, "image not exists")
		return
	}
	if req.FuzzerID != 0 && !service.IsStorageItemExistsByID(req.FuzzerID) {
		utils.BadRequestWithMsg(c, "fuzzer not exists")
		return
	}
	if req.CorpusID != 0 && !service.IsStorageItemExistsByID(req.CorpusID) {
		utils.BadRequestWithMsg(c, "corpus not exists")
		return
	}
	if req.TargetID != 0 && !service.IsStorageItemExistsByID(req.TargetID) {
		utils.BadRequestWithMsg(c, "target not exists")
		return
	}

	err = service.CreateTask(&task)
	if err != nil {
		utils.DBError(c)
		return
	}
	var Err error
	defer func() {
		if Err != nil {
			service.DeleteTask(task.ID)
		}
	}()
	if req.Environments != nil {
		err = service.InsertEnvironments(task.ID, req.Environments)
		if err != nil {
			Err = err
			utils.DBError(c)
			return
		}
	}
	if req.Arguments != nil {
		err = service.InsertArguments(task.ID, req.Arguments)
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
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '202':
	//     description: start task success
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	task, err := getTask(c)
	if err != nil {
		return
	}
	if task.Status != models.TaskCreated {
		utils.BadRequestWithMsg(c, "wrong status")
		return
	}
	var Err error
	if !service.IsStorageItemExistsByID(task.FuzzerID) {
		utils.BadRequestWithMsg(c, "you should upload fuzzer first")
		return
	}
	if !service.IsStorageItemExistsByID(task.TargetID) {
		utils.BadRequestWithMsg(c, "you should upload target first")
		return
	}
	if !service.IsStorageItemExistsByID(task.CorpusID) {
		utils.BadRequestWithMsg(c, "you should upload corpus first")
		return
	}

	err = kubernetes.CreateServiceByTaskID(task.ID)
	if err != nil {
		utils.InternalErrorWithMsg(c, "create service failed")
		return
	}

	defer func() {
		if Err != nil {
			kubernetes.DeleteServiceByTaskID(task.ID)
		}
	}()
	tempImage, err := service.GetImageByID(task.ImageID)
	if err != nil {
		Err = err
		utils.DBError(c)
		return
	}
	if tempImage == nil {
		Err = errors.New("image not exists")
		utils.BadRequestWithMsg(c, "image not exists")
		return
	}
	var image *appsv1.Deployment
	if !tempImage.IsDeployment {
		image, err = kubernetes.GenerateDeployment(task.ID, tempImage.Content, 1)
		if err != nil {
			Err = err
			utils.InternalErrorWithMsg(c, "generate image failed")
			return
		}
	} else {
		image, err = kubernetes.GenerateDeploymentByYaml(tempImage.Content, task.ID)
		if err != nil {
			Err = err
			utils.BadRequestWithMsg(c, err.Error())
			return
		}
	}
	err = kubernetes.CreateDeploy(image)
	if err != nil {
		Err = err
		utils.InternalErrorWithMsg(c, "create image failed: "+err.Error())
		return
	}
	if err = models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Status", models.TaskStarted).Error; err != nil {
		kubernetes.DeleteDeployByTaskID(task.ID)
		utils.DBError(c)
		Err = err
		return
	}
	if err = models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("StatusUpdateAt", time.Now().Unix()).Error; err != nil {
		kubernetes.DeleteDeployByTaskID(task.ID)
		utils.DBError(c)
		Err = err
		return
	}
	c.String(http.StatusAccepted, "")
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
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '202':
	//     description: stop task success
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	task, err := getTask(c)
	if err != nil {
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
	err = kubernetes.DeleteContainerByTaskID(task.ID)
	if err != nil {
		utils.InternalErrorWithMsg(c, "kubernetes delete error")
		return
	}
	c.String(http.StatusAccepted, "")
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
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	task, err := getTask(c)
	if err != nil {
		return
	}
	var req TaskUpdateReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}

	if task.Status != models.TaskCreated {
		utils.BadRequestWithMsg(c, "you can change task after it started")
		return
	}
	updateData := map[string]interface{}{}

	if req.ImageID != 0 {
		if service.IsObjectExistsByID(&models.Image{}, req.ImageID) {
			updateData["ImageID"] = req.ImageID
		} else {
			utils.BadRequestWithMsg(c, "image id not exist")
			return
		}
	}
	ids := []uint64{req.FuzzerID, req.CorpusID, req.TargetID}
	types := []string{models.Fuzzer, models.Corpus, models.Target}
	modelsField := []string{"FuzzerID", "CorpusID", "TargetID"}
	for i, _ := range ids {
		if ids[i] != 0 {
			storageItem, err := service.GetStorageItemByID(ids[i])
			if err != nil {
				utils.DBError(c)
				return
			}
			if storageItem == nil {
				utils.BadRequestWithMsg(c, types[i]+" not exist")
				return
			}
			if storageItem.Type != types[i] {
				utils.BadRequestWithMsg(c, "wrong type")
				return
			}
			updateData[modelsField[i]] = ids[i]
		}
	}
	if req.Time != 0 {
		updateData["Time"] = req.Time
	}
	if req.FuzzCycleTime != 0 {
		updateData["FuzzCycleTime"] = req.FuzzCycleTime
	}
	if req.Arguments != nil {
		if err = service.DeleteObjectsByTaskID(&models.TaskArgument{}, task.ID); err != nil {
			utils.DBError(c)
			return
		}
		if err = service.InsertArguments(task.ID, req.Arguments); err != nil {
			utils.DBError(c)
			return
		}
	}
	if req.Environments != nil {
		if err = service.DeleteObjectsByTaskID(&models.TaskEnvironment{}, task.ID); err != nil {
			utils.DBError(c)
			return
		}
		if err = service.InsertEnvironments(task.ID, req.Environments); err != nil {
			utils.DBError(c)
			return
		}
	}
	if err = service.UpdateTask(task, updateData); err != nil {
		utils.DBError(c)
		return
	}

	c.String(http.StatusCreated, "")
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
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	task, err := getTask(c)
	if err != nil {
		return
	}
	if task.UserID != uint64(c.GetInt64("id")) && !c.GetBool("isAdmin") {
		utils.Forbidden(c)
		return
	}
	if task.Status != models.TaskCreated && task.Status != models.TaskStopped && task.Status != models.TaskError {
		utils.BadRequestWithMsg(c, "task is still running")
		return
	}
	if err := service.DeleteTask(task.ID); err != nil {
		utils.DBError(c)
		return
	}
	kubernetes.DeleteContainerByTaskID(task.ID)
	c.String(http.StatusNoContent, "")
}
