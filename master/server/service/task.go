package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/jinzhu/gorm"
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

func DeleteObjectsByTaskID(obj interface{}, taskid uint64) error {
	return models.DB.Where("task_id = ?", taskid).Delete(obj).Error
}

func CreateTask(task *models.Task) error {
	return insertObject(task)
}

func UpdateTask(task *models.Task, data map[string]interface{}) error {
	return UpdateObject(task, data)
}
