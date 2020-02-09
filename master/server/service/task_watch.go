package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"time"
)

func isDeployReady(deploy *appsv1.Deployment) bool {
	logger.Logger.Debug("deployment status", "status", deploy.Status)
	return deploy.Status.AvailableReplicas >= 1
}

func isPodFailed(pod *corev1.Pod) bool {
	for _, status := range pod.Status.ContainerStatuses {
		logger.Logger.Debug("container status", "status", status.State.Waiting)
		if status.State.Waiting == nil {
			continue
		}
		if status.State.Waiting.Reason == "ImagePullBackOff" {
			return true
		}
	}
	return false
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
				deploy, ok := newObj.(*appsv1.Deployment)
				if ok && isDeployReady(deploy) {
					go initDeployTask(deploy)
				}
			},
		},
	)
	go controller.Run(deployWatchChan)
}

func watchPod() {
	watchlist := cache.NewListWatchFromClient(ClientSet.CoreV1().RESTClient(), "pods", config.KubernetesConf.Namespace, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&corev1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod, ok := newObj.(*corev1.Pod)
				if ok && isPodFailed(pod) {
					go deleteTaskByPod(pod)
				}
			},
		},
	)
	go controller.Run(podWatchChan)
}

func deleteTaskByPod(pod *corev1.Pod) {
	taskID, err := getPodTaskID(pod)
	if err != nil {
		return
	}
	logger.Logger.Debug("deployment failed to start")
	SetTaskError(taskID, "failed to create deployment")
}
