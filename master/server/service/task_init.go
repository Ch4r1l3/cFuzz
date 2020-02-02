package service

import (
	"encoding/json"
	botmodels "github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"time"
)

func initDeployTask(deploy *appsv1.Deployment) {
	var taskID uint64
	var task models.Task
	var err, Err error
	taskID, err = getDeployTaskID(deploy)
	if err != nil {
		return
	}
	defer func() {
		if Err != nil {
			logger.Logger.Error("Error exit init", "reason", Err.Error())
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskError)
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix())
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("ErrorMsg", "DB Error")
			DeleteDeployByTaskID(taskID)
			DeleteServiceByTaskID(taskID)
		}
	}()
	if err = models.GetObjectByID(&task, uint64(taskID)); err != nil {
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
		DeleteDeployByTaskID(taskID)
		DeleteServiceByTaskID(taskID)
		return
	}
	//test if bot is not up, retry 3 times
	for i := 0; i < 3; i++ {
		_, err = requestProxyGet(uint64(taskID), []string{"storage_item"})
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
		if err = models.GetObjectByID(&storageItem, ids[i]); err != nil {
			Err = errors.Wrap(err, "DB Error")
			return
		}
		if storageItem.ExistsInImage {
			uploadFuzzerPostData := map[string]interface{}{
				"type":          types[i],
				"existsInImage": true,
				"path":          storageItem.Path,
			}
			result, err := requestProxyPost(uint64(taskID), []string{"storage_item", "exist"}, uploadFuzzerPostData)
			if err != nil {
				Err = errors.Wrap(err, "upload fuzzer Error")
				return
			}
			var resp clientStorageItemPostResp
			if err := json.Unmarshal(result, &resp); err != nil {
				Err = errors.Wrap(err, "json decode fuzzer resp Error")
				return
			}
			botids = append(botids, resp.ID)
		} else {
			form := map[string]string{
				"type":    types[i],
				"relPath": storageItem.RelPath,
			}
			result, err := requestProxyPostWithFile(uint64(taskID), []string{"storage_item"}, form, storageItem.Path)
			if err != nil {
				Err = errors.Wrap(err, "upload fuzzer Error")
				return
			}
			var resp clientStorageItemPostResp
			if err := json.Unmarshal(result, &resp); err != nil {
				Err = errors.Wrap(err, "json decode fuzzer resp Error")
				return
			}
			botids = append(botids, resp.ID)
		}
	}

	taskArguments, err := models.GetArguments(task.ID)
	if err != nil {
		Err = errors.Wrap(err, "DB Error")
		return
	}
	taskEnvironments, err := models.GetEnvironments(task.ID)
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
	//result, err = requestProxyPost(task.ID, []string{"task"}, postData)
	result, err := requestProxyPostRaw(task.ID, []string{"task"}, postData)
	logger.Logger.Debug("create task", "result", string(result))
	if err != nil {
		Err = errors.Wrap(err, "create task error")
		return
	}

	//start fuzz target
	result, err = requestProxyPost(task.ID, []string{"task", "start"}, struct{}{})
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
