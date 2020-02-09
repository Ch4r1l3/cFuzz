package service

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"strconv"
	"time"
)

func GetAllDeploys() ([]appsv1.Deployment, error) {
	deploys, err := ClientSet.AppsV1().Deployments(config.KubernetesConf.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return deploys.Items, nil
}

func GetAllServices() ([]apiv1.Service, error) {
	services, err := ClientSet.CoreV1().Services(config.KubernetesConf.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return services.Items, nil
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

func GenerateDeployment(taskID uint64, taskName string, image string, replicasNum int32) (*appsv1.Deployment, error) {
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
							Name:  taskName,
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

func CreateServiceByTaskID(taskID uint64) error {
	_, err := ClientSet.CoreV1().Services(config.KubernetesConf.Namespace).Create(&apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf(ServiceNameFmt, taskID),
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Selector: map[string]string{
				LabelName: fmt.Sprintf("%d", taskID),
			},
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Protocol:   apiv1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	})
	return err
}

func DeleteServiceByTaskID(taskID uint64) error {
	return ClientSet.CoreV1().Services(config.KubernetesConf.Namespace).Delete(fmt.Sprintf(ServiceNameFmt, taskID), &metav1.DeleteOptions{})
}

func DeleteContainerByTaskID(taskID uint64) error {
	err1 := DeleteDeployByTaskID(taskID)
	err2 := DeleteServiceByTaskID(taskID)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func getDeployTaskID(deploy *appsv1.Deployment) (uint64, error) {
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

func getServiceTaskID(service *apiv1.Service) (uint64, error) {
	if service.ObjectMeta.Labels == nil {
		return 0, errors.New("label not exists")
	}
	if _, ok := service.ObjectMeta.Labels[LabelName]; !ok {
		return 0, errors.New("taskid not exists in label")
	}
	taskID, err := strconv.ParseInt(service.ObjectMeta.Labels[LabelName], 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(taskID), nil
}

func SetTaskError(taskID uint64, errorMsg string) {
	models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("Status", models.TaskError)
	models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("StatusUpdateAt", time.Now().Unix())
	models.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("ErrorMsg", errorMsg)
	DeleteContainerByTaskID(taskID)
}
