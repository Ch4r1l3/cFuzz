package service

import (
	"encoding/json"
	botmodels "github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"
)

func checkTasks() {
	for {
		var tasks []models.Task
		if err := models.GetObjects(&tasks); err != nil {
			logger.Logger.Error("checkTasks", "error", err.Error())
		} else {
			for _, task := range tasks {
				if task.Status == models.TaskStarted && task.StatusUpdateAt+config.KubernetesConf.MaxStartTime < time.Now().Unix() {
					go func() {
						logger.Logger.Error("deployment start too long", "start", task.StatusUpdateAt, "now", time.Now().Unix())
						models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Status", models.TaskError)
						models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("StatusUpdateAt", time.Now().Unix())
						models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("ErrorMsg", "deployment start waste to much time")
						DeleteDeployByTaskID(task.ID)
						DeleteServiceByTaskID(task.ID)
					}()
				} else if task.Status == models.TaskRunning && time.Now().Unix()-task.StartedAt > int64(task.Time) {
					go func() {
						models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Status", models.TaskStopped)
						models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("StatusUpdateAt", time.Now().Unix())
						DeleteDeployByTaskID(task.ID)
						DeleteServiceByTaskID(task.ID)
					}()
				} else if task.Status == models.TaskRunning {
					if activeRoutineNum[task.ID] == nil {
						activeRoutineNum[task.ID] = new(int32)
					}
					if *activeRoutineNum[task.ID] == 0 {
						atomic.AddInt32(activeRoutineNum[task.ID], 1)
						ants.Submit(func() {
							checkSingleTask(task.ID)
						})
					}
				}

			}
		}
		<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime) * time.Second)
	}

}

func initCrashesMap() error {
	crashesMap = make(map[uint64]map[uint64]bool)
	var crashes []models.TaskCrash
	if err := models.GetObjects(&crashes); err != nil {
		return err
	}

	for _, crash := range crashes {
		if _, ok := crashesMap[crash.TaskID]; !ok {
			crashesMap[crash.TaskID] = make(map[uint64]bool)
		}
		crashesMap[crash.TaskID][crash.BotCrashID] = true
	}
	return nil
}

func checkSingleTask(taskID uint64) {
	defer func() {
		tempNumPtr := activeRoutineNum[taskID]
		if tempNumPtr != nil {
			atomic.AddInt32(tempNumPtr, -1)
		}
	}()
	var Err error
	defer func() {
		if Err != nil {
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskError)
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix())
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("ErrorMsg", Err.Error())
			DeleteDeployByTaskID(taskID)
			DeleteServiceByTaskID(taskID)
		}
	}()
	if err := models.DB.Model(&models.Task{}).
		Where("id = ?", taskID).Update("Status", models.TaskRunning).Error; err != nil {
		logger.Logger.Error("DB error", "reason", err.Error())
		return
	}
	if err := models.DB.Model(&models.Task{}).
		Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix()).Error; err != nil {
		logger.Logger.Error("DB error", "reason", err.Error())
		return
	}
	var task models.Task
	if err := models.GetObjectByID(&task, taskID); err != nil {
		logger.Logger.Error("DB error", "reason", err.Error())
		return
	}
	if task.Status != models.TaskRunning {
		return
	}

	//get deployment check if it is running
	deploys, err := GetDeployByTaskID(taskID)
	if len(deploys) == 0 {
		logger.Logger.Error("get deployment empty")
		Err = errors.New("get deployment empty")
		return
	}
	if err != nil {
		logger.Logger.Error("get deployment error", "reason", err.Error())
		Err = err
		return
	}

	//get bot task status
	result, err := requestProxyGet(taskID, []string{"task"})
	if err != nil {
		logger.Logger.Error("get task error", "reason", err.Error())
		Err = err
		return
	}
	var clientTask clientTaskGetResp
	if err := json.Unmarshal(result, &clientTask); err != nil {
		logger.Logger.Error("get task error", "reason", err.Error())
		Err = err
		return
	}
	if clientTask.Status != botmodels.TaskRunning {
		logger.Logger.Debug("client status is not running", "status", clientTask.Status)
		if clientTask.Status == botmodels.TaskError {
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskError)
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("ErrorMsg", "client error exit")
		} else {
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskStopped)
		}
		models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix())
		DeleteServiceByTaskID(taskID)
		DeleteDeployByTaskID(taskID)
		return
	} else {
		//get bot task crashes
		result, err = requestProxyGet(taskID, []string{"task", "crash"})
		if err != nil {
			logger.Logger.Error("client get crash failed", "reason", err.Error())
			Err = err
			return
		}
		var crashes []clientCrashGetResp
		if err = json.Unmarshal(result, &crashes); err != nil {
			logger.Logger.Error("client  crash json decode fail", "reason", err.Error())
			return
		}
		logger.Logger.Debug("task", "crashes", crashes)
		for _, crash := range crashes {
			crashesPath := filepath.Join(config.ServerConf.CrashesPath, strconv.Itoa(int(taskID)))
			if _, ok := crashesMap[taskID]; !ok {
				crashesMap[taskID] = make(map[uint64]bool)
				os.MkdirAll(crashesPath, os.ModePerm)
			}
			if !crashesMap[taskID][crash.ID] {
				crashesMap[taskID][crash.ID] = true
				savePath, err := requestProxySaveFile(taskID, []string{"task", "crash", strconv.Itoa(int(crash.ID))}, crashesPath)
				if err != nil {
					logger.Logger.Error("request save file error", "reason", err.Error())
				}
				taskCrash := models.TaskCrash{
					BotCrashID:    crash.ID,
					TaskID:        taskID,
					Path:          savePath,
					ReproduceAble: crash.ReproduceAble,
				}
				if err := models.DB.Create(&taskCrash).Error; err != nil {
					Err = err
					return
				}
			}
		}

		//get bot task result
		result, err = requestProxyGet(taskID, []string{"task", "result"})
		if err != nil {
			logger.Logger.Error("client get result failed", "reason", err.Error())
			Err = err
			return
		}
		if len(result) > 10 {
			var fuzzResult clientResultGetResp
			if err = json.Unmarshal(result, &fuzzResult); err != nil {
				logger.Logger.Error("client fuzz result json decode fail", "len", len(result), "content", result, "reason", err.Error())
				return
			}
			lastFuzzResult, err := models.GetLastestFuzzResultByTaskID(taskID)
			if err != nil {
				logger.Logger.Error("get fuzz result from db error", "reason", err.Error())
				Err = err
				return
			}
			logger.Logger.Debug("last fuzz result", "result", lastFuzzResult)
			if lastFuzzResult.UpdateAt < fuzzResult.UpdateAt {
				taskFuzzResult := models.TaskFuzzResult{
					Command:      fuzzResult.Command,
					TimeExecuted: fuzzResult.TimeExecuted,
					TaskID:       taskID,
					UpdateAt:     fuzzResult.UpdateAt,
				}
				if err = models.DB.Create(&taskFuzzResult).Error; err != nil {
					Err = err
					return
				}
				if err = models.InsertTaskFuzzResultStat(taskFuzzResult.ID, fuzzResult.Stats); err != nil {
					Err = err
					return
				}
			}
		}
	}
}
