package models

import (
	"github.com/jinzhu/gorm"
)

type Task struct {
	ID             uint64 `gorm:"primary_key" json:"id"`
	Name           string `json:"name"`
	Image          string `json:"image"`
	DeploymentID   uint64 `json:"deploymentid"`
	Time           uint64 `json:"time"`
	FuzzCycleTime  uint64 `json:"fuzzCycleTime"`
	FuzzerID       uint64 `json:"fuzzerid"`
	CorpusID       uint64 `json:"corpusid"`
	TargetID       uint64 `json:"targetid"`
	Status         string `json:"status"`
	ErrorMsg       string `json:"errorMsg"`
	StatusUpdateAt int64  `json:"-"`
	StartedAt      int64  `json:"startedAt"`
}

const (
	TaskRunning      = "TaskRunning"
	TaskStarted      = "TaskStarted"
	TaskCreated      = "TaskCreated"
	TaskInitializing = "TaskInitializing"
	TaskStopped      = "TaskStopped"
	TaskError        = "TaskError"
)

type TaskEnvironment struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	Value  string `json:"value"`
}

func DeleteTask(taskid uint64) error {
	var err error
	if err = DeleteObjectsByTaskID(&TaskEnvironment{}, taskid); err != nil {
		return err
	}
	if err = DeleteObjectsByTaskID(&TaskArgument{}, taskid); err != nil {
		return err
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
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
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

func DeleteObjectsByTaskID(obj interface{}, taskid uint64) error {
	return DB.Where("task_id = ?", taskid).Delete(obj).Error
}

func GetObjectsByTaskID(obj interface{}, taskid uint64) error {
	return DB.Where("task_id = ?", taskid).Find(obj).Error
}

func GetObjectByTaskIDAndID(obj interface{}, taskid uint64, id uint64) error {
	return DB.Where("task_id = ? AND id = ?", taskid, id).Error
}

type TaskCrash struct {
	ID            uint64 `gorm:"primary_key;auto_increment:false" json:"id"`
	TaskID        uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	Path          string `json:"-"`
	ReproduceAble bool   `json:"reproduceAble"`
}

type TaskFuzzResult struct {
	ID           uint64 `gorm:"primary_key" json:"id"`
	Command      string `json:"command"`
	TimeExecuted int    `json:"timeExecuted"`
	TaskID       uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	UpdateAt     int64  `json:"updateAt"`
}

type TaskFuzzResultStat struct {
	Key              string `json:"key"`
	Value            string `json:"value"`
	TaskFuzzResultID uint64 `json:"taskid" sql:"type:bigint REFERENCES task_fuzz_result(id) ON DELETE CASCADE"`
}

func InsertTaskFuzzResultStat(id uint64, stats map[string]string) error {
	for k, v := range stats {
		taskFuzzResultStat := TaskFuzzResultStat{
			TaskFuzzResultID: id,
			Key:              k,
			Value:            v,
		}
		if err := DB.Create(&taskFuzzResultStat).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetTaskFuzzResultStat(id uint64) (map[string]string, error) {
	var taskFuzzResultStat []TaskFuzzResultStat
	if err := DB.Where("task_fuzz_result_id = ?", id).Find(&taskFuzzResultStat).Error; err != nil {
		return nil, err
	}
	stats := make(map[string]string)
	for _, v := range taskFuzzResultStat {
		stats[v.Key] = v.Value
	}
	return stats, nil
}

func GetLastestFuzzResultByTaskID(taskid uint64) (*TaskFuzzResult, error) {
	var taskFuzzResult TaskFuzzResult
	if err := DB.Where("task_id = ?", taskid).Order("update_at desc").First(&taskFuzzResult).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &taskFuzzResult, nil
		}
		return nil, err
	}
	return &taskFuzzResult, nil
}
