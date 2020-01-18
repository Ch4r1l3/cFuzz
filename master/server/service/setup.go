package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

var ClientSet kubernetes.Interface
var controlChan map[uint64]chan struct{}
var crashesMap map[uint64]map[uint64]bool
var deployWatchChan chan struct{}

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
	controlChan = make(map[uint64]chan struct{})
	crashesMap = make(map[uint64]map[uint64]bool)
	deployWatchChan = make(chan struct{})
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

	setupNamespace()
	watchDeploy()
	watchPod()
	go handleTasks()
}
