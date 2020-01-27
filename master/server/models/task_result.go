package models

import (
	"github.com/jinzhu/gorm"
)

// swagger:model
type TaskFuzzResult struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`
	// example: /afl/afl-fuzz -i xx -o xx ./test
	Command string `json:"command"`
	// example: 60
	TimeExecuted int `json:"timeExecuted"`
	// example: 1
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	// example: 1579996805
	UpdateAt int64 `json:"updateAt"`
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
