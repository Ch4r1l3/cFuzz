package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/jinzhu/gorm"
)

func InsertTaskFuzzResultStat(id uint64, stats map[string]string) error {
	for k, v := range stats {
		taskFuzzResultStat := models.TaskFuzzResultStat{
			TaskFuzzResultID: id,
			Key:              k,
			Value:            v,
		}
		if err := models.DB.Create(&taskFuzzResultStat).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetTaskFuzzResultStat(id uint64) (map[string]string, error) {
	var taskFuzzResultStat []models.TaskFuzzResultStat
	if err := models.DB.Where("task_fuzz_result_id = ?", id).Find(&taskFuzzResultStat).Error; err != nil {
		return nil, err
	}
	stats := make(map[string]string)
	for _, v := range taskFuzzResultStat {
		stats[v.Key] = v.Value
	}
	return stats, nil
}

func GetLastestFuzzResultByTaskID(taskid uint64) (*models.TaskFuzzResult, error) {
	var taskFuzzResult models.TaskFuzzResult
	if err := models.DB.Where("task_id = ?", taskid).Order("update_at desc").First(&taskFuzzResult).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &taskFuzzResult, nil
}
