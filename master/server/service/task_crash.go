package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
)

func GetTaskCrashes() ([]models.TaskCrash, error) {
	var crashes []models.TaskCrash
	if err := getObjects(&crashes); err != nil {
		return nil, err
	}
	return crashes, nil
}
