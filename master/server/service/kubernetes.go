package service

import (
	"errors"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
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
	if deployment.Spec.Selector.MatchLabels == nil {
		deployment.Spec.Selector.MatchLabels = make(map[string]string)
	}
	deployment.Spec.Selector.MatchLabels[LabelName] = fmt.Sprintf("%d", taskID)
	if deployment.Spec.Template.ObjectMeta.Labels == nil {
		deployment.Spec.Template.ObjectMeta.Labels = make(map[string]string)
	}
	deployment.Spec.Template.ObjectMeta.Labels[LabelName] = fmt.Sprintf("%d", taskID)
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

func TestProxy() {
	restClient := ClientSet.CoreV1().RESTClient()
	request := restClient.Get().Namespace("default").Resource("services").Name("hostname-service").Suffix("proxy", "50x.html")
	fmt.Println(request.URL())
	result := request.Do()
	var statusCode int
	result.StatusCode(&statusCode)
	fmt.Println(statusCode)
	bytes, _ := result.Raw()
	fmt.Println(string(bytes))
	fmt.Println(result.Error())

}

func handleTasks() {
	for {
		var tasks []models.Task
		if err := models.GetObjects(&tasks); err != nil {
			fmt.Println(err)
		} else {
			for _, task := range tasks {
				if task.Running {
					result, err, statusCode := requestProxyGet(task.ID, []string{"fuzzer"})
					if err == nil {
						fmt.Println("handleTask")
						fmt.Println(statusCode)
						fmt.Println(result)
						fmt.Println("handleTask")
					}
				}
			}
		}
		<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime) * time.Second)
	}

}

func setupNewPod() {
	watchlist := cache.NewListWatchFromClient(ClientSet.AppsV1().RESTClient(), "deployments", config.KubernetesConf.Namespace,
		fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&appsv1.Deployment{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				if deploy, ok := newObj.(*appsv1.Deployment); ok {
					for _, v := range deploy.Status.Conditions {
						fmt.Printf("condition: %s ;", v.Type)
					}
					fmt.Println("")
					if deploy.ObjectMeta.Labels == nil {
						return
					}
					if _, ok := deploy.ObjectMeta.Labels[LabelName]; !ok {
						return
					}
					if len(deploy.Status.Conditions) >= 1 && deploy.Status.Conditions[0].Type == appsv1.DeploymentAvailable {
						taskID, err := strconv.ParseInt(deploy.ObjectMeta.Labels[LabelName], 10, 64)
						if err != nil {
							return
						}
						var task models.Task
						if err = models.GetObjectByID(&task, taskID); err != nil {
							return
						}
						var fuzzer models.Fuzzer
						if err = models.GetObjectByID(&fuzzer, task.FuzzerID); err != nil {
							return
						}
						result, err, statusCode := requestProxyGet(uint64(taskID), []string{"fuzzer"})
						if err != nil {
							return
						}
						fmt.Println(result)
						form := map[string]string{
							"name", fuzzer.Name,
						}
						result, err, statusCode = requestProxyPostWithFile(taskID, []string{"fuzzer"}, form, fuzzer.Path)
						if err != nil {
							return
						}

					}
				}
			},
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
}
