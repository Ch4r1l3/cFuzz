package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// swagger:model
type TaskFuzzResult struct {
	// example: /afl/afl-fuzz -i xxx -o xxx ./test
	Command string `json:"command"`
	// example: 60
	TimeExecuted int `json:"timeExecuted"`
	// example: 1579996805
	UpdateAt int64 `json:"updateAt"`
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
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	resultStats := []TaskFuzzResultStat{}
	if err := DB.Find(&resultStats).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	stats := make(map[string]string)
	for _, v := range resultStats {
		stats[v.Key] = v.Value
	}
	return &result, stats, nil

}
