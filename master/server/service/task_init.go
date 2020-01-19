package service

import (
	"encoding/json"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
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
	var fuzzer models.Fuzzer
	if err = models.GetObjectByID(&fuzzer, task.FuzzerID); err != nil {
		Err = errors.Wrap(err, "DB Error")
		return
	}
	//test if bot is not up, retry 3 times
	for i := 0; i < 3; i++ {
		_, err = requestProxyGet(uint64(taskID), []string{"fuzzer"})
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
	//upload fuzzer to bot
	form := map[string]string{
		"name": fuzzer.Name,
	}
	result, err := requestProxyPostWithFile(uint64(taskID), []string{"fuzzer"}, form, fuzzer.Path)
	if err != nil {
		Err = errors.Wrap(err, "upload fuzzer Error")
		return
	}
	var clientFuzzer clientFuzzerPostResp
	if err := json.Unmarshal(result, &clientFuzzer); err != nil {
		Err = errors.Wrap(err, "json decode fuzzer resp Error")
		return
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
		"fuzzerID":      clientFuzzer.ID,
		"maxTime":       task.Time,
		"fuzzCycleTime": task.FuzzCycleTime,
		"arguments":     taskArguments,
		"enviroments":   taskEnvironments,
	}

	//create task on bot
	//result, err = requestProxyPost(task.ID, []string{"task"}, postData)
	result, err = requestProxyPostRaw(task.ID, []string{"task"}, postData)
	logger.Logger.Debug("create task", "result", string(result))
	if err != nil {
		Err = errors.Wrap(err, "create task error")
		return
	}
	var taskTarget []models.TaskTarget
	if err = models.GetObjectsByTaskID(&taskTarget, uint64(taskID)); err != nil || len(taskTarget) == 0 {
		if err != nil {
			Err = errors.Wrap(err, "get task target error")
		} else {
			Err = errors.New("get task target error target empty")
		}
		return
	}
	var taskCorpus []models.TaskCorpus
	if err = models.GetObjectsByTaskID(&taskCorpus, uint64(taskID)); err != nil || len(taskCorpus) == 0 {
		if err != nil {
			Err = errors.Wrap(err, "get task corpus error")
		} else {
			Err = errors.New("get task corpus error corpus empty")
		}
		return
	}

	//upload target and corpus to bot
	result, err = requestProxyPostWithFile(uint64(taskID), []string{"task", "target"}, form, taskTarget[0].Path)
	if err != nil {
		Err = errors.Wrap(err, "upload target error")
		return
	}
	result, err = requestProxyPostWithFile(uint64(taskID), []string{"task", "corpus"}, form, taskCorpus[0].Path)
	if err != nil {
		Err = errors.Wrap(err, "upload corpus error")
		return
	}

	//start fuzz target
	putData := map[string]interface{}{
		"status": "TASK_RUNNING",
	}
	result, err = requestProxyPut(task.ID, []string{"task"}, putData)
	if err != nil {
		Err = errors.Wrap(err, "start bot fuzz error")
		return
	}

	v, ok := controlChan[task.ID]
	if ok {
		v <- struct{}{}
	} else {
		controlChan[task.ID] = make(chan struct{})
		go func() {
			<-time.After(time.Duration(task.Time) * time.Second)
			logger.Logger.Debug("time end!")
			controlChan[task.ID] <- struct{}{}
			DeleteServiceByTaskID(taskID)
			DeleteDeployByTaskID(taskID)
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskStopped)
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix())
			<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime) * time.Second)
			delete(controlChan, task.ID)
			delete(crashesMap, task.ID)
		}()
	}
	checkSingleTask(taskID)
}
