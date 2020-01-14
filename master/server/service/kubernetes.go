package service

import (
	"encoding/json"
	"errors"
	"fmt"
	botmodels "github.com/Ch4r1l3/cFuzz/bot/server/models"
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
		count := make(map[uint64]int)
		crashesMap := make(map[uint64]map[uint64]bool)
		if err := models.GetObjects(&tasks); err != nil {
			fmt.Println(err)
		} else {
			for _, task := range tasks {
				if task.Running {
					result, err, statusCode := requestProxyGet(task.ID, []string{"fuzzer"})
					if err != nil || statusCode >= 300 {
						v, ok := count[task.ID]
						if ok {
							count[task.ID] = v + 1
						} else {
							count[task.ID] = 1
						}
						if v >= config.KubernetesConf.MaxClientRetryNum {
							delete(count, task.ID)
							if _, ok = crashesMap[task.ID]; ok {
								delete(crashesMap, task.ID)
							}
							models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Running", false)
						}
						return
					}
					fmt.Printf("fuzzer status Code: %d\n", statusCode)
					result, err, statusCode = requestProxyGet(task.ID, []string{"task"})
					var clientTask clientTaskGetResp
					if err := json.Unmarshal(result, &clientTask); err != nil {
						fmt.Println(err.Error())
						return
					}
					fmt.Println(string(result))
					if clientTask.Status != botmodels.TASK_RUNNING {
						fmt.Println("client status is not running, status: " + clientTask.Status + ".")
						models.DB.Model(&models.Task{}).Where("id = ?", task.ID).Update("Running", false)
						DeleteServiceByTaskID(task.ID)
						DeleteDeployByTaskID(task.ID)
						return
					} else {
						result, err, statusCode = requestProxyGet(task.ID, []string{"task", "crash"})
						if err != nil {
							return
						}
						var crashes []clientCrashGetResp
						if err = json.Unmarshal(result, &crashes); err != nil {
							fmt.Println(err.Error())
							return
						}
						fmt.Println("handle task")
						fmt.Println(statusCode)
						fmt.Println(crashes)
						for _, crash := range crashes {
							if _, ok := crashesMap[task.ID]; !ok {
								crashesMap[task.ID] = make(map[uint64]bool)
							}
							if _, ok := crashesMap[task.ID][crash.ID]; !ok {
								crashesMap[task.ID][crash.ID] = true
								result, err, statusCode = requestProxyGet(task.ID, []string{"task", "crash", strconv.Itoa(int(crash.ID))})
								fmt.Println(result)
							}
						}
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
						<-time.After(time.Duration(3) * time.Second)
						if err != nil {
							return
						}
						var task models.Task
						if err = models.GetObjectByID(&task, uint64(taskID)); err != nil {
							fmt.Println(err.Error())
							return
						}
						var fuzzer models.Fuzzer
						if err = models.GetObjectByID(&fuzzer, task.FuzzerID); err != nil {
							fmt.Println(err.Error())
							return
						}
						result, err, statusCode := requestProxyGet(uint64(taskID), []string{"fuzzer"})
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						form := map[string]string{
							"name": fuzzer.Name,
						}
						result, err, statusCode = requestProxyPostWithFile(uint64(taskID), []string{"fuzzer"}, form, fuzzer.Path)
						if err != nil {
							return
						}
						if statusCode > 300 {
							return
						}
						var clientFuzzer clientFuzzerPostResp
						if err := json.Unmarshal(result, &clientFuzzer); err != nil {
							return
						}
						if !ok {
							fmt.Println("cannot change to struct")
							return
						}
						taskArguments, err := models.GetArguments(task.ID)
						if err != nil {
							return
						}
						taskEnvironments, err := models.GetEnvironments(task.ID)
						if err != nil {
							return
						}
						postData := map[string]interface{}{
							"fuzzerID":      clientFuzzer.ID,
							"maxTime":       task.Time,
							"fuzzCycleTime": task.FuzzCycleTime,
							"arguments":     taskArguments,
							"enviroments":   taskEnvironments,
						}
						//result, err, statusCode = requestProxyPost(task.ID, []string{"task"}, postData)
						fmt.Printf("cycleTime %d\n", task.FuzzCycleTime)
						result, err = requestProxyPostRaw(task.ID, []string{"task"}, postData)
						fmt.Printf("create task result: %s\n", string(result))
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						var taskTarget []models.TaskTarget
						if err = models.GetObjectsByTaskID(&taskTarget, uint64(taskID)); err != nil || len(taskTarget) == 0 {
							return
						}
						var taskCorpus []models.TaskCorpus
						if err = models.GetObjectsByTaskID(&taskCorpus, uint64(taskID)); err != nil || len(taskCorpus) == 0 {
							return
						}
						result, err, statusCode = requestProxyPostWithFile(uint64(taskID), []string{"task/target"}, form, taskTarget[0].Path)
						if err != nil {
							return
						}
						result, err, statusCode = requestProxyPostWithFile(uint64(taskID), []string{"task/corpus"}, form, taskCorpus[0].Path)
						if err != nil {
							return
						}
						putData := map[string]interface{}{
							"status": "TASK_RUNNING",
						}
						result, err, statusCode = requestProxyPut(task.ID, []string{"task"}, putData)
					}
				}
			},
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
}
