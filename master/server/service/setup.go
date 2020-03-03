package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service/kubernetes"
	"log"
)

var crashesMap map[uint64]map[uint64]bool
var activeRoutineNum map[uint64]*int32

func Setup() {
	activeRoutineNum = make(map[uint64]*int32)
	kubernetes.Setup()

	initKubernetesCleanup()
	initCrashesMap()
	kubernetes.WatchDeploy(initDeployTask)
	kubernetes.WatchPod(func(taskID uint64) {
		SetTaskError(taskID, "failed to create deployment")
	})
	go checkTasks()
}

func initKubernetesCleanup() {
	var tasks []models.Task
	if err := getObjects(&tasks); err != nil {
		logger.Logger.Error("checkTasks", "error", err.Error())
		return
	}
	tempRecord := make(map[uint64]bool)
	for _, task := range tasks {
		if task.Status == models.TaskStarted || task.Status == models.TaskInitializing || (config.KubernetesConf.InitCleanup && task.Status == models.TaskRunning) {
			SetTaskError(task.ID, "server stopped")
		}
		if task.Status == models.TaskRunning {
			tempRecord[task.ID] = true
		}
	}
	deploys, err := kubernetes.GetAllDeploys()
	if err != nil {
		log.Fatal(err)
	}
	for _, deploy := range deploys {
		taskID, err := kubernetes.GetDeployTaskID(&deploy)
		if err != nil {
			continue
		}
		// clean all those deployment that is not running
		if !tempRecord[taskID] {
			kubernetes.DeleteDeployByTaskID(taskID)
		}
	}
	services, err := kubernetes.GetAllServices()
	if err != nil {
		log.Fatal(err)
	}
	for _, service := range services {
		taskID, err := kubernetes.GetServiceTaskID(&service)
		if err != nil {
			continue
		}
		// clean all those service that is not running
		if !tempRecord[taskID] {
			kubernetes.DeleteServiceByTaskID(taskID)
		}
	}
}
