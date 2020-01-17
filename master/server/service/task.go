package service

import (
	"encoding/json"
	"fmt"
	botmodels "github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"strconv"
	"time"
)

type clientFuzzerPostResp struct {
	ID   uint64 `json:"id" binding:"required"`
	Name string `json:"string" binding:"required"`
}

type clientTaskGetResp struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type clientCrashGetResp struct {
	ID            uint64 `json:"id" binding:"required"`
	ReproduceAble bool   `json:"reproduceAble" binding:"required"`
}

func isDeployReady(deploy *appsv1.Deployment) bool {
	logger.Logger.Debug("deployment status", "status", deploy.Status)
	logger.Logger.Debug("deployment OwnerReferences", "OwnerReferences", deploy.ObjectMeta.OwnerReferences)
	return deploy.Status.AvailableReplicas >= 1
}

func watchDeploy() {
	watchlist := cache.NewListWatchFromClient(ClientSet.AppsV1().RESTClient(), "deployments", config.KubernetesConf.Namespace,
		fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&appsv1.Deployment{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				if deploy, ok := newObj.(*appsv1.Deployment); ok && isDeployReady(deploy) {
					go initDeployTask(deploy)
				}
			},
		},
	)
	go controller.Run(deployWatchChan)
}

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
			logger.Logger.Warn("Error exit init", "reason", Err.Error())
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskError)
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
			<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime) * time.Second)
			delete(controlChan, task.ID)
			delete(crashesMap, task.ID)
		}()
	}
	handleSingleTask(taskID)
}

func handleTasks() {
	for {
		var tasks []models.Task
		if err := models.GetObjects(&tasks); err != nil {
			fmt.Println(err)
		} else {
			for _, task := range tasks {
				if task.Status == models.TaskInitializing {
				}
			}
		}
		<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime) * time.Second)
	}

}

func handleSingleTask(taskID uint64) {
	var Err error
	errRetryNum := 0
	goReturn := func() bool {
		errRetryNum += 1
		return errRetryNum >= config.KubernetesConf.MaxClientRetryNum
	}
	defer func() {
		if Err != nil {
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskError)
			models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("ErrorMsg", Err.Error())
		}
	}()
	if err := models.DB.Model(&models.Task{}).
		Where("id = ?", taskID).Update("Status", models.TaskRunning).Error; err != nil {
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
					if _, ok := crashesMap[taskID]; !ok {
						crashesMap[taskID] = make(map[uint64]bool)
					}
					if _, ok := crashesMap[taskID][crash.ID]; !ok {
						crashesMap[taskID][crash.ID] = true
						result, err = requestProxyGet(taskID, []string{"task", "crash", strconv.Itoa(int(crash.ID))})
						logger.Logger.Debug("get task crash", "result", result)
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
