package service

import (
	"encoding/json"
	"fmt"
	botmodels "github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func checkTasks() {
	for {
		var tasks []models.Task
		if err := models.GetObjects(&tasks); err != nil {
			fmt.Println(err)
		} else {
			for _, task := range tasks {
				if task.Status == models.TaskStarted && task.StatusUpdateAt+config.KubernetesConf.MaxStartTime < time.Now().Unix() {
					logger.Logger.Error("deployment start too long", "start", task.StatusUpdateAt, "now", time.Now().Unix())
					models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Status", models.TaskError)
					models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("StatusUpdateAt", time.Now().Unix())
					models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("ErrorMsg", "deployment start waste to much time")
					DeleteDeployByTaskID(task.ID)
					DeleteServiceByTaskID(task.ID)
				}
			}
		}
		<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime) * time.Second)
	}

}

func checkSingleTask(taskID uint64) {
	var Err error
	errRetryNum := 0
	var lastUpdateTime int64
	lastUpdateTime = 0
	goReturn := func() bool {
		errRetryNum += 1
		return errRetryNum >= config.KubernetesConf.MaxClientRetryNum
	}
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
		Err = err
		return
	}
	if err := models.DB.Model(&models.Task{}).
		Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix()).Error; err != nil {
		logger.Logger.Error("DB error", "reason", err.Error())
		Err = err
		return
	}
	for {
		select {
		case <-time.After(time.Duration(config.KubernetesConf.CheckTaskTime) * time.Second):
			var task models.Task
			if err := models.GetObjectByID(&task, taskID); err != nil {
				logger.Logger.Error("DB error", "reason", err.Error())
				if goReturn() {
					Err = err
					return
				} else {
					continue
				}
			}
			if task.Status != models.TaskRunning {
				if task.Status == models.TaskInitializing {
					<-controlChan[taskID]
				}
				if goReturn() {
					return
				} else {
					continue
				}
			}
			deploys, err := GetDeployByTaskID(taskID)
			if len(deploys) == 0 {
				logger.Logger.Error("get deployment empty")
				if goReturn() {
					Err = errors.New("get deployment empty")
					return
				} else {
					continue
				}
			}
			if err != nil {
				logger.Logger.Error("get deployment error", "reason", err.Error())
				if goReturn() {
					Err = err
					return
				} else {
					continue
				}
			}
			result, err := requestProxyGet(taskID, []string{"task"})
			if err != nil {
				logger.Logger.Error("get task error", "reason", err.Error())
				if goReturn() {
					Err = err
					return
				} else {
					continue
				}
			}
			var clientTask clientTaskGetResp
			if err := json.Unmarshal(result, &clientTask); err != nil {
				logger.Logger.Error("get task error", "reason", err.Error())
				if goReturn() {
					Err = err
					return
				} else {
					continue
				}
			}
			if clientTask.Status != botmodels.TASK_RUNNING {
				logger.Logger.Debug("client status is not running", "status", clientTask.Status)
				if clientTask.Status == botmodels.TASK_ERROR {
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
				result, err = requestProxyGet(taskID, []string{"task", "crash"})
				if err != nil {
					logger.Logger.Error("client get crash failed", "reason", err.Error())
					if goReturn() {
						Err = err
						return
					} else {
						continue
					}
				}
				var crashes []clientCrashGetResp
				if err = json.Unmarshal(result, &crashes); err != nil {
					logger.Logger.Error("client  crash json decode fail", "reason", err.Error())
					if goReturn() {
						return
					} else {
						continue
					}
				}
				logger.Logger.Debug("task", "crashes", crashes)
				for _, crash := range crashes {
					crashesPath := filepath.Join(config.ServerConf.CrashesPath, strconv.Itoa(int(taskID)))
					if _, ok := crashesMap[taskID]; !ok {
						crashesMap[taskID] = make(map[uint64]bool)
						os.MkdirAll(crashesPath, os.ModePerm)
					}
					if _, ok := crashesMap[taskID][crash.ID]; !ok {
						crashesMap[taskID][crash.ID] = true
						savePath, err := requestProxySaveFile(taskID, []string{"task", "crash", strconv.Itoa(int(crash.ID))}, crashesPath)
						if err != nil {
							logger.Logger.Error("request save file error", "reason", err.Error())
						}
						taskCrash := models.TaskCrash{
							TaskID: taskID,
							Path:   savePath,
						}
						if err := models.DB.Create(&taskCrash).Error; err != nil {
							if goReturn() {
								Err = err
								return
							} else {
								continue
							}
						}
					}
				}
				result, err = requestProxyGet(taskID, []string{"task", "result"})
				if err != nil {
					logger.Logger.Error("client get result failed", "reason", err.Error())
					if goReturn() {
						Err = err
						return
					} else {
						continue
					}
				}
				if len(result) > 10 {
					var fuzzResult clientResultGetResp
					if err = json.Unmarshal(result, &fuzzResult); err != nil {
						logger.Logger.Error("client fuzz result json decode fail", "len", len(result), "content", result, "reason", err.Error())
						if goReturn() {
							return
						} else {
							continue
						}
					}
					if lastUpdateTime != fuzzResult.UpdateAt {
						lastUpdateTime = fuzzResult.UpdateAt
						taskFuzzResult := models.TaskFuzzResult{
							Command:      fuzzResult.Command,
							TimeExecuted: fuzzResult.TimeExecuted,
							TaskID:       taskID,
							UpdateAt:     fuzzResult.UpdateAt,
						}
						if err = models.DB.Create(&taskFuzzResult).Error; err != nil {
							if goReturn() {
								Err = err
								return
							} else {
								continue
							}
						}
						if err = models.InsertTaskFuzzResultStat(taskFuzzResult.ID, fuzzResult.Stats); err != nil {
							if goReturn() {
								Err = err
								return
							} else {
								continue
							}
						}
					}
				}
			}
			errRetryNum = 0

		case <-controlChan[taskID]:
			logger.Logger.Debug("recv stop signal")
			return
		}
	}

}
