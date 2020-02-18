package models

import (
	"github.com/jinzhu/gorm"
)

// swagger:model
type Task struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: test
	Name string `json:"name" sql:"type:varchar(255) NOT NULL UNIQUE"`

	// example: 1
	ImageID uint64 `json:"imageID" sql:"type:integer REFERENCES image(id)"`

	// example: 60
	Time uint64 `json:"time"`

	// example: 60
	FuzzCycleTime uint64 `json:"fuzzCycleTime"`

	// example: 1
	FuzzerID uint64 `json:"fuzzerID" sql:"type:integer REFERENCES storage_item(id)"`

	// example: 2
	CorpusID uint64 `json:"corpusID" sql:"type:integer REFERENCES storage_item(id)"`

	// example: 3
	TargetID uint64 `json:"targetID" sql:"type:integer REFERENCES storage_item(id)"`

	// example: TaskRunning
	Status string `json:"status"`

	// example: pull image error
	ErrorMsg string `json:"errorMsg"`

	// example: 1579996805
	StatusUpdateAt int64 `json:"-"`

	// example: 1579996805
	StartedAt int64 `json:"startedAt"`

	// example: 1
	UserID uint64 `json:"userID" sql:"type:integer REFERENCES user(id) ON DELETE CASCADE"`
}

const (
	TaskRunning      = "TaskRunning"
	TaskStarted      = "TaskStarted"
	TaskCreated      = "TaskCreated"
	TaskInitializing = "TaskInitializing"
	TaskStopped      = "TaskStopped"
	TaskError        = "TaskError"
)

func (t *Task) IsRunning() bool {
	return t.Status == TaskStarted || t.Status == TaskInitializing || t.Status == TaskRunning
}

type TaskEnvironment struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:integer REFERENCES task(id) ON DELETE CASCADE"`
	Value  string `json:"value"`
}

func DeleteTask(taskid uint64) error {
	var err error
	// delete crashes
	var crashes []TaskCrash
	if err = GetObjectsByTaskID(&crashes, taskid); err != nil {
		return err
	}
	for _, c := range crashes {
		c.Delete()
	}
	if err = DeleteObjectByID(&Task{}, taskid); err != nil {
		return err
	}
	return nil
}

func InsertEnvironments(taskid uint64, environments []string) error {
	for _, v := range environments {
		taskEnvironment := TaskEnvironment{
			TaskID: taskid,
			Value:  v,
		}
		if err := DB.Create(&taskEnvironment).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetEnvironments(taskid uint64) ([]string, error) {
	var taskEnvironments []TaskEnvironment
	if err := DB.Where("task_id = ?", taskid).Find(&taskEnvironments).Error; err != nil {
		return nil, err
	}
	environments := []string{}
	for _, v := range taskEnvironments {
		environments = append(environments, v.Value)
	}
	return environments, nil
}

type TaskArgument struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:integer REFERENCES task(id) ON DELETE CASCADE"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func InsertArguments(taskid uint64, arguments map[string]string) error {
	for k, v := range arguments {
		taskArgument := TaskArgument{
			TaskID: taskid,
			Key:    k,
			Value:  v,
		}
		if err := DB.Create(&taskArgument).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetArguments(taskid uint64) (map[string]string, error) {
	var taskArguments []TaskArgument
	if err := DB.Where("task_id = ?", taskid).Find(&taskArguments).Error; err != nil {
		return nil, err
	}
	arguments := make(map[string]string)
	for _, v := range taskArguments {
		arguments[v.Key] = v.Value
	}
	return arguments, nil
}

func GetObjectsByTaskID(obj interface{}, taskid uint64) error {
	return DB.Where("task_id = ?", taskid).Find(obj).Error
}

func GetCountByTaskID(obj interface{}, taskid uint64) (int, error) {
	var count int
	err := DB.Model(obj).Where("task_id = ?", taskid).Count(&count).Error
	return count, err
}

func GetObjectsByTaskIDPagination(objs interface{}, taskid uint64, offset int, limit int) (int, error) {
	return GetObjectCombinCustom(objs, offset, limit, "", []string{"task_id = ?"}, []interface{}{taskid})
}

func GetTaskByID(id uint64) (*Task, error) {
	var task Task
	if err := DB.Where("id = ?", id).First(&task).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func DeleteObjectsByTaskID(obj interface{}, taskid uint64) error {
	return DB.Where("task_id = ?", taskid).Delete(obj).Error
}
