package models

import (
	"time"
)

type Task struct {
	CorpusDir     string `json:"corpusDir"`
	TargetDir     string `json:"targetDir"`
	TargetPath    string `json:"targetPath"`
	Status        string `json:"status"`
	FuzzerID      uint64 `json:"fuzzerID"`
	FuzzCycleTime uint64 `json:"fuzzCycleTime"`
	MaxTime       int    `json:"maxTime"`
}

const (
	TASK_CREATED = "TASK_CREATED"
	TASK_RUNNING = "TASK_RUNNING"
	TASK_STOPPED = "TASK_STOPPED"
	TASK_ERROR   = "TASK_ERROR"
)

func GetTask() (*Task, error) {
	var task Task
	if err := DB.First(&task).Error; err != nil {
		return &Task{}, err
	}
	return &task, nil
}

type TaskArgument struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func InsertArguments(arguments map[string]string) error {
	for k, v := range arguments {
		ta := TaskArgument{
			Key:   k,
			Value: v,
		}
		result := DB.Create(&ta)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func GetArguments() (map[string]string, error) {
	arguments := make(map[string]string)
	taskArguments := []TaskArgument{}
	if err := DB.Find(&taskArguments).Error; err != nil {
		return nil, err
	}
	for _, v := range taskArguments {
		arguments[v.Key] = v.Value
	}
	return arguments, nil
}

type TaskEnvironment struct {
	Value string `json:"value"`
}

func InsertEnvironments(environments []string) error {
	for _, v := range environments {
		te := TaskEnvironment{
			Value: v,
		}
		result := DB.Create(&te)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func GetEnvironments() ([]string, error) {
	environments := []string{}
	taskEnvironments := []TaskEnvironment{}
	if err := DB.Find(&taskEnvironments).Error; err != nil {
		return environments, err
	}
	for _, v := range taskEnvironments {
		environments = append(environments, v.Value)
	}
	return environments, nil
}

type TaskCrash struct {
	ID            uint64 `gorm:"primary_key";json:"id"`
	Path          string `json:"path"`
	ReproduceAble bool   `json:"reproduceAble"`
}

func GetCrashes() ([]TaskCrash, error) {
	taskCrashes := []TaskCrash{}
	if err := DB.Find(&taskCrashes).Error; err != nil {
		return nil, err
	}
	return taskCrashes, nil
}

func GetCrashByID(id uint64) (*TaskCrash, error) {
	var crash TaskCrash
	if err := DB.Where("id = ?", id).First(&crash).Error; err != nil {
		return nil, err
	}
	return &crash, nil
}

func CreateCrash(path string, reproduceAble bool) error {
	taskCrash := TaskCrash{
		Path:          path,
		ReproduceAble: reproduceAble,
	}
	return DB.Save(&taskCrash).Error
}

type TaskFuzzResult struct {
	Command      string `json:"command"`
	TimeExecuted int    `json:"timeExecuted"`
	UpdateAt     int64  `json:"updateAt"`
}

type TaskFuzzResultStat struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func CreateFuzzResult(command []string, stats map[string]string, timeExecuted int) error {
	DB.Delete(&TaskFuzzResult{})
	DB.Delete(&TaskFuzzResultStat{})
	tcommand := ""
	for _, v := range command {
		tcommand += v + " "
	}
	fuzzResult := TaskFuzzResult{
		Command:      tcommand,
		TimeExecuted: timeExecuted,
		UpdateAt:     time.Now().Unix(),
	}
	if err := DB.Save(&fuzzResult).Error; err != nil {
		return err
	}
	for k, v := range stats {
		stat := TaskFuzzResultStat{
			Key:   k,
			Value: v,
		}
		if err := DB.Save(&stat).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetFuzzResult() (*TaskFuzzResult, map[string]string, error) {
	var result TaskFuzzResult
	if err := DB.First(&result).Error; err != nil {
		return nil, nil, err
	}
	resultStats := []TaskFuzzResultStat{}
	if err := DB.Find(&resultStats).Error; err != nil {
		return nil, nil, err
	}
	stats := make(map[string]string)
	for _, v := range resultStats {
		stats[v.Key] = v.Value
	}
	return &result, stats, nil

}
