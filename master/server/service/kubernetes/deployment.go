package kubernetes

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/logger"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"strconv"
	"time"
)

type functype func(uint64)

func GetAllDeploys() ([]appsv1.Deployment, error) {
	deploys, err := ClientSet.AppsV1().Deployments(config.KubernetesConf.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return deploys.Items, nil
}

func GetDeployByTaskID(taskID uint64) ([]appsv1.Deployment, error) {
	labelSelector := fmt.Sprintf(LabelFmt, taskID)
	deploys, err := ClientSet.AppsV1().Deployments(config.KubernetesConf.Namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}
	return deploys.Items, nil
}

func CreateDeploy(deployment *appsv1.Deployment) error {
	_, err := ClientSet.AppsV1().Deployments(config.KubernetesConf.Namespace).Create(deployment)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDeployByTaskID(taskID uint64) error {
	return ClientSet.AppsV1().Deployments(config.KubernetesConf.Namespace).Delete(fmt.Sprintf(DeployNameFmt, taskID), &metav1.DeleteOptions{})
}

func GenerateDeploymentByYaml(content string, taskID uint64) (*appsv1.Deployment, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(content), nil, nil)
	if err != nil {
		return nil, err
	}
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, errors.New("can not parse it to deployment")
	}
	if deployment.ObjectMeta.Labels == nil {
		deployment.ObjectMeta.Labels = make(map[string]string)
	}
	deployment.ObjectMeta.Labels[LabelName] = fmt.Sprintf("%d", taskID)
	if deployment.Spec.Selector == nil {
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: map[string]string{},
		}
	}
	if deployment.Spec.Selector.MatchLabels == nil {
		deployment.Spec.Selector.MatchLabels = make(map[string]string)
	}
	deployment.Spec.Selector.MatchLabels[LabelName] = fmt.Sprintf("%d", taskID)
	if deployment.Spec.Template.ObjectMeta.Labels == nil {
		deployment.Spec.Template.ObjectMeta.Labels = make(map[string]string)
	}
	deployment.Spec.Template.ObjectMeta.Labels[LabelName] = fmt.Sprintf("%d", taskID)
	deployment.ObjectMeta.Name = fmt.Sprintf(DeployNameFmt, taskID)
	return deployment, nil
}

func GenerateDeployment(taskID uint64, image string, replicasNum int32) (*appsv1.Deployment, error) {
	if replicasNum <= 0 {
		return nil, errors.New("replicas number should large than 0")
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf(DeployNameFmt, taskID),
			Labels: map[string]string{
				LabelName: fmt.Sprintf("%d", taskID),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicasNum,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					LabelName: fmt.Sprintf("%d", taskID),
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						LabelName: fmt.Sprintf("%d", taskID),
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  fmt.Sprintf(DeployNameFmt, taskID),
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment, nil
}

func GetDeployTaskID(deploy *appsv1.Deployment) (uint64, error) {
	if deploy.ObjectMeta.Labels == nil {
		return 0, errors.New("label not exists")
	}
	if _, ok := deploy.ObjectMeta.Labels[LabelName]; !ok {
		return 0, errors.New("taskid not exists in label")
	}
	taskID, err := strconv.ParseInt(deploy.ObjectMeta.Labels[LabelName], 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(taskID), nil
}

func isDeployReady(deploy *appsv1.Deployment) bool {
	logger.Logger.Debug("deployment status", "status", deploy.Status)
	return deploy.Status.AvailableReplicas >= 1
}

func WatchDeploy(callback functype) {
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
					taskID, err := GetDeployTaskID(deploy)
					if err != nil {
						return
					}
					go callback(taskID)
				}
			},
		},
	)
	go controller.Run(deployWatchChan)
}
