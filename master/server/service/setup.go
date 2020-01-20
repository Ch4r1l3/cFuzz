package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

var ClientSet kubernetes.Interface
var crashesMap map[uint64]map[uint64]bool
var activeRoutineNum map[uint64]*int32
var deployWatchChan chan struct{}
var podWatchChan chan struct{}

func setupNamespace() {
	if config.KubernetesConf.Namespace == "" {
		log.Fatal("namespace can not be empty")
	}
	_, err := ClientSet.CoreV1().Namespaces().Get(config.KubernetesConf.Namespace, metav1.GetOptions{})
	if err != nil {
		_, err = ClientSet.CoreV1().Namespaces().Create(&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: config.KubernetesConf.Namespace,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

const (
	LabelName      = "taskID"
	LabelFmt       = "taskID=%d"
	DeployNameFmt  = "task%d"
	ServiceNameFmt = "task%d-service"
)

func Setup() {
	deployWatchChan = make(chan struct{})
	podWatchChan = make(chan struct{})
	activeRoutineNum = make(map[uint64]*int32)
	var kubeConfig *rest.Config
	var err error
	if config.KubernetesConf.ConfigPath != "" {
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.KubernetesConf.ConfigPath)

	} else if home := homedir.HomeDir(); home != "" {
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	} else {
		kubeConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		log.Fatal(err)
	}
	ClientSet, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	initKubernetesCleanup()
	initCrashesMap()
	setupNamespace()
	watchDeploy()
	watchPod()
	go checkTasks()
}

func initKubernetesCleanup() {
	var tasks []models.Task
	if err := models.GetObjects(&tasks); err != nil {
		logger.Logger.Error("checkTasks", "error", err.Error())
		return
	}
	tempRecord := make(map[uint64]bool)
	for _, task := range tasks {
		if task.Status == models.TaskStarted || task.Status == models.TaskInitializing || (config.KubernetesConf.InitCleanup && task.Status == models.TaskRunning) {
			models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Status", models.TaskError)
			models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("StatusUpdateAt", time.Now().Unix())
			models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("ErrorMsg", "server stopped")
		}
		if task.Status == models.TaskRunning {
			tempRecord[task.ID] = true
		}
	}
	deploys, err := GetAllDeploys()
	if err != nil {
		log.Fatal(err)
	}
	for _, deploy := range deploys {
		taskID, err := getDeployTaskID(&deploy)
		if err != nil {
			continue
		}
		// clean all those deployment that is not running
		if !tempRecord[taskID] {
			DeleteDeployByTaskID(taskID)
		}
	}
	services, err := GetAllServices()
	if err != nil {
		log.Fatal(err)
	}
	for _, service := range services {
		taskID, err := getServiceTaskID(&service)
		if err != nil {
			continue
		}
		// clean all those service that is not running
		if !tempRecord[taskID] {
			DeleteServiceByTaskID(taskID)
		}
	}

}
