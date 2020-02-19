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
	var task models.Task
	var err, Err error
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
	if err = GetObjectByID(&task, taskID); err != nil {
		Err = err
		return
	}
	if task.Status == models.TaskStarted || task.Status == models.TaskRunning {
		if err := models.DB.Model(&models.Task{}).
			Where("id = ?", taskID).Update("Status", models.TaskInitializing).Error; err != nil {
			Err = errors.Wrap(err, "DB Error")
			return
		}
		if err := models.DB.Model(&models.Task{}).
			Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix()).Error; err != nil {
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
		var storageItem models.StorageItem
		if err = GetObjectByID(&storageItem, ids[i]); err != nil {
			Err = errors.Wrap(err, "DB Error")
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

	if err := models.DB.Model(&models.Task{}).
		Where("id = ?", taskID).Update("Status", models.TaskRunning).Error; err != nil {
		logger.Logger.Error("DB error", "reason", err.Error())
		Err = errors.Wrap(err, "DB error")
		return
	}
	if err := models.DB.Model(&models.Task{}).
		Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix()).Error; err != nil {
		logger.Logger.Error("DB error", "reason", err.Error())
		Err = errors.Wrap(err, "DB error")
		return
	}
	if task.StartedAt == 0 {
		if err := models.DB.Model(&models.Task{}).
			Where("id = ?", taskID).Update("StartedAt", time.Now().Unix()).Error; err != nil {
			logger.Logger.Error("DB error", "reason", err.Error())
			Err = errors.Wrap(err, "DB error")
			return
		}
	}
}
