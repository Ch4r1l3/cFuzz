package service

import (
	botmodels "github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service/kubernetes"
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
		var err error
		if tasks, err = GetTasks(); err != nil {
			logger.Logger.Error("checkTasks", "error", err.Error())
		} else {
			for _, task := range tasks {
				//logger.Logger.Debug("checkTasks")
				if task.Status == models.TaskStarted && task.StatusUpdateAt+config.KubernetesConf.MaxStartTime < time.Now().Unix() {
					go func() {
						logger.Logger.Error("deployment start too long", "start", task.StatusUpdateAt, "now", time.Now().Unix())
						SetTaskError(task.ID, "deployment start waste to much time")
					}()
				} else if task.Status == models.TaskRunning && time.Now().Unix()-task.StartedAt > int64(task.Time) {
					go func() {
						SetTaskStopped(task.ID)
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
	crashes, err := GetTaskCrashes()
	if err != nil {
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
			//check current task status first
			tempTask, err := GetTaskByID(uint64(taskID))
			if tempTask == nil || err != nil {
				return
			}
			if tempTask.Status != models.TaskError && tempTask.Status != models.TaskStopped {
				SetTaskError(taskID, Err.Error())
			}
		}
	}()
	task, err := GetTaskByID(taskID)
	if err != nil {
		logger.Logger.Error("DB error", "reason", err.Error())
		return
	}
	if task == nil {
		logger.Logger.Error("Task not exists", "id", taskID)
		return
	}
	if task.Status != models.TaskRunning {
		return
	}

	//get deployment check if it is running
	deploys, err := kubernetes.GetDeployByTaskID(taskID)
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
	clientStatus, clientErrorMsg, err := kubernetes.GetTask(taskID)
	if err != nil {
		logger.Logger.Error("get task error", "reason", err.Error())
		Err = err
		return
	}
	if clientStatus != botmodels.TaskRunning {
		logger.Logger.Debug("client status is not running", "status", clientStatus)
		if clientStatus == botmodels.TaskError {
			SetTaskError(taskID, "client error exit: "+clientErrorMsg)
		} else {
			SetTaskStopped(taskID)
		}
		return
	} else {
		//get bot task crashes
		crashes, err := kubernetes.GetCrashes(taskID)
		if err != nil {
			logger.Logger.Error("client get crashes fail", "reason", err.Error())
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
				savePath, err := kubernetes.DownloadCrash(taskID, crash.ID, crashesPath)
				if err != nil {
					logger.Logger.Error("request save file error", "reason", err.Error())
					continue
				}
				taskCrash := models.TaskCrash{
					BotCrashID:    crash.ID,
					TaskID:        taskID,
					Path:          savePath,
					ReproduceAble: crash.ReproduceAble,
					FileName:      crash.FileName,
				}
				if err := models.DB.Create(&taskCrash).Error; err != nil {
					Err = err
					return
				}
			}
		}

		//get bot task result
		fuzzResult, err := kubernetes.GetResult(taskID)
		if err != nil {
			logger.Logger.Error("client get result fail", "reason", err.Error())
			return
		}
		if fuzzResult != nil {
			lastFuzzResult, err := GetLastestFuzzResultByTaskID(taskID)
			if err != nil {
				logger.Logger.Error("get fuzz result from db error", "reason", err.Error())
				Err = err
				return
			}
			logger.Logger.Debug("last fuzz result", "result", lastFuzzResult)
			if lastFuzzResult == nil || lastFuzzResult.UpdateAt < fuzzResult.UpdateAt {
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
				if err = InsertTaskFuzzResultStat(taskFuzzResult.ID, fuzzResult.Stats); err != nil {
					Err = err
					return
				}
			}
		}
	}
}
