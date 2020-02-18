package kubernetes

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"strconv"
	"time"
)

func getPodTaskID(pod *apiv1.Pod) (uint64, error) {
	if pod.ObjectMeta.Labels == nil {
		return 0, errors.New("label not exists")
	}
	if _, ok := pod.ObjectMeta.Labels[LabelName]; !ok {
		return 0, errors.New("taskid not exists in label")
	}
	taskID, err := strconv.ParseInt(pod.ObjectMeta.Labels[LabelName], 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(taskID), nil
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

func WatchPod(callback functype) {
	watchlist := cache.NewListWatchFromClient(ClientSet.CoreV1().RESTClient(), "pods", config.KubernetesConf.Namespace, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&corev1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod, ok := newObj.(*corev1.Pod)
				if ok && isPodFailed(pod) {
					taskID, err := getPodTaskID(pod)
					if err != nil {
						return
					}
					logger.Logger.Debug("deployment failed to start")
					go callback(taskID)
				}
			},
		},
	)
	go controller.Run(podWatchChan)
}
