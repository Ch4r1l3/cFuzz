package service

import (
	botmodels "github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service/kubernetes"
	"github.com/pkg/errors"
	"time"
)

func initDeployTask(taskID uint64) {
	var Err error
	defer func() {
		if Err != nil {
			logger.Logger.Error("Error exit init", "reason", Err.Error())
			// check current task status first
			tempTask, err := GetTaskByID(taskID)
			if tempTask == nil || err != nil {
				return
			}
			if tempTask.Status != models.TaskError && tempTask.Status != models.TaskStopped {
				SetTaskError(taskID, "init task error: "+Err.Error())
			}
		}
	}()
	task, err := GetTaskByID(taskID)
	if err != nil || task == nil {
		Err = err
		return
	}
	if task.Status == models.TaskStarted || task.Status == models.TaskRunning {
		if err := UpdateTaskStatus(taskID, models.TaskInitializing); err != nil {
			Err = errors.Wrap(err, "DB Error")
			return
		}
	} else {
		kubernetes.DeleteDeployByTaskID(taskID)
		kubernetes.DeleteServiceByTaskID(taskID)
		return
	}
	//test if bot is not up, retry 3 times
	for i := 0; i < 3; i++ {
		_, err = kubernetes.GetStorageItems(task.ID)
		if err != nil {
			if i == 2 {
				Err = errors.Wrap(err, "service Error")
				return
			}
		} else {
			break
		}
		<-time.After(time.Second)
	}
	//upload fuzzer, corpus, target
	ids := []uint64{task.FuzzerID, task.CorpusID, task.TargetID}
	types := []string{botmodels.Fuzzer, botmodels.Corpus, botmodels.Target}
	botids := []uint64{}
	for i, _ := range ids {
		storageItem, err := GetStorageItemByID(ids[i])
		if err != nil {
			Err = errors.Wrap(err, "DB Error")
			return
		}
		if storageItem == nil {
			Err = errors.New("storageItem not exists")
			return
		}
		tid, err := kubernetes.CreateStorageItem(task.ID, storageItem.ExistsInImage, types[i], storageItem.Path, storageItem.RelPath)
		if err != nil {
			Err = errors.Wrap(err, "upload storageItem failed")
			return
		}
		botids = append(botids, tid)
	}

	taskArguments, err := GetArguments(task.ID)
	if err != nil {
		Err = errors.Wrap(err, "DB Error")
		return
	}
	taskEnvironments, err := GetEnvironments(task.ID)
	if err != nil {
		Err = errors.Wrap(err, "DB Error")
		return
	}
	postData := map[string]interface{}{
		"fuzzerID":      botids[0],
		"corpusID":      botids[1],
		"targetID":      botids[2],
		"maxTime":       task.Time,
		"fuzzCycleTime": task.FuzzCycleTime,
		"arguments":     taskArguments,
		"enviroments":   taskEnvironments,
	}

	//create task on bot
	err = kubernetes.CreateTask(task.ID, postData)
	if err != nil {
		Err = errors.Wrap(err, "create task error")
		return
	}
	//start fuzz target
	err = kubernetes.StartTask(task.ID)
	if err != nil {
		Err = errors.Wrap(err, "start bot fuzz error")
		return
	}
	//update task status
	if err := UpdateTaskStatus(taskID, models.TaskRunning); err != nil {
		Err = errors.Wrap(err, "DB Error")
		logger.Logger.Error("DB error", "reason", err.Error())
		return
	}
	if task.StartedAt == 0 {
		if err := UpdateTaskField(taskID, "StartedAt", time.Now().Unix()); err != nil {
			logger.Logger.Error("DB error", "reason", err.Error())
			Err = errors.Wrap(err, "DB error")
			return
		}
	}
}
