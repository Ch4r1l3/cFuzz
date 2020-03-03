package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service/kubernetes"
	"github.com/imroc/req"
	"github.com/jinzhu/gorm"
	"time"
)

func DeleteTask(taskid uint64) error {
	var err error
	// delete crashes
	var crashes []models.TaskCrash
	if err = GetObjectsByTaskID(&crashes, taskid); err != nil {
		return err
	}
	for _, c := range crashes {
		c.Delete()
	}
	if err = DeleteObjectByID(&models.Task{}, taskid); err != nil {
		return err
	}
	return nil
}

func InsertEnvironments(taskid uint64, environments []string) error {
	for _, v := range environments {
		taskEnvironment := models.TaskEnvironment{
			TaskID: taskid,
			Value:  v,
		}
		if err := models.DB.Create(&taskEnvironment).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetEnvironments(taskid uint64) ([]string, error) {
	var taskEnvironments []models.TaskEnvironment
	if err := models.DB.Where("task_id = ?", taskid).Find(&taskEnvironments).Error; err != nil {
		return nil, err
	}
	environments := []string{}
	for _, v := range taskEnvironments {
		environments = append(environments, v.Value)
	}
	return environments, nil
}

func InsertArguments(taskid uint64, arguments map[string]string) error {
	for k, v := range arguments {
		taskArgument := models.TaskArgument{
			TaskID: taskid,
			Key:    k,
			Value:  v,
		}
		if err := models.DB.Create(&taskArgument).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetArguments(taskid uint64) (map[string]string, error) {
	var taskArguments []models.TaskArgument
	if err := models.DB.Where("task_id = ?", taskid).Find(&taskArguments).Error; err != nil {
		return nil, err
	}
	arguments := make(map[string]string)
	for _, v := range taskArguments {
		arguments[v.Key] = v.Value
	}
	return arguments, nil
}

func GetObjectsByTaskID(obj interface{}, taskid uint64) error {
	return models.DB.Where("task_id = ?", taskid).Find(obj).Error
}

func GetCountByTaskID(obj interface{}, taskid uint64) (int, error) {
	var count int
	err := models.DB.Model(obj).Where("task_id = ?", taskid).Count(&count).Error
	return count, err
}

func GetObjectsByTaskIDPagination(objs interface{}, taskid uint64, offset int, limit int) (int, error) {
	return getObjectCombinCustom(objs, offset, limit, "", []string{"task_id = ?"}, []interface{}{taskid})
}

func GetTaskByID(id uint64) (*models.Task, error) {
	var task models.Task
	if err := GetObjectByID(&task, id); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := getObjects(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByUserID(userID uint64) ([]models.Task, error) {
	var tasks []models.Task
	_, err := getObjectCombinCustom(&tasks, -1, -1, "", []string{"user_id = ?"}, []interface{}{userID})
	return tasks, err
}

func DeleteObjectsByTaskID(obj interface{}, taskid uint64) error {
	return models.DB.Where("task_id = ?", taskid).Delete(obj).Error
}

func CreateTask(task *models.Task) error {
	return insertObject(task)
}

func UpdateTask(task *models.Task, data map[string]interface{}) error {
	return UpdateObject(task, data)
}

func UpdateTaskField(id uint64, name string, value interface{}) error {
	return UpdateObjectField(&models.Task{}, id, name, value)
}

func UpdateTaskStatus(id uint64, Status string) error {
	if err := UpdateTaskField(id, "StatusUpdateAt", time.Now().Unix()); err != nil {
		return err
	}
	return UpdateTaskField(id, "Status", Status)
}

func SetTaskError(id uint64, errorMsg string) error {
	UpdateTaskStatus(id, models.TaskError)
	UpdateTaskField(id, "ErrorMsg", errorMsg)
	ReqCallbackUrl(id)
	return kubernetes.DeleteContainerByTaskID(id)
}

func SetTaskStopped(id uint64) error {
	UpdateTaskStatus(id, models.TaskStopped)
	ReqCallbackUrl(id)
	return kubernetes.DeleteContainerByTaskID(id)
}

func ReqCallbackUrl(id uint64) {
	task, err := GetTaskByID(id)
	if err != nil || task == nil {
		return
	}
	r := req.New()
	req.SetTimeout(10 * time.Second)
	param := req.Param{
		"taskID": id,
	}
	for i := 0; i < 3; i++ {
		resp, err := r.Post(task.CallbackUrl, param)
		if err == nil && resp.Response().StatusCode < 300 {
			break
		}
	}
}
