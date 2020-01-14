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

func handleSingleTask(taskID uint64) {
	if err := models.DB.Model(&models.Task{}).
		Where("id = ?", taskID).Update("Status", models.TaskRunning).Error; err != nil {
		logger.Logger.Error("DB error", "reason", err.Error())
		return
	}
	for {
		select {
		case <-time.After(time.Duration(config.KubernetesConf.CheckTaskTime) * time.Second):
			var task models.Task
			if err := models.GetObjectByID(&task, taskID); err != nil {
				logger.Logger.Error("DB error", "reason", err.Error())
				return
			}
			if task.Status != models.TaskRunning {

			}
			result, err := requestProxyGet(taskID, []string{"task"})
			var clientTask clientTaskGetResp
			if err := json.Unmarshal(result, &clientTask); err != nil {
				logger.Logger.Debug("get task error", "reason", err.Error())
				return
			}
			if clientTask.Status != botmodels.TASK_RUNNING {
				logger.Logger.Debug("client status is not running", "status", clientTask.Status)
				DeleteServiceByTaskID(taskID)
				DeleteDeployByTaskID(taskID)
				return
			} else {
				result, err = requestProxyGet(taskID, []string{"task", "crash"})
				if err != nil {
					return
				}
				var crashes []clientCrashGetResp
				if err = json.Unmarshal(result, &crashes); err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("handle task")
				fmt.Println(crashes)
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

		case <-controlChan[taskID]:
			return
		}
	}

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
			Err = errors.Wrap(err, "service Error")
			return
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
	}
	go func() {
		<-time.After(time.Duration(task.Time) * time.Second)
		controlChan[task.ID] <- struct{}{}
		<-time.After(10)
		delete(controlChan, task.ID)
		delete(crashesMap, task.ID)
	}()
	handleSingleTask(taskID)
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
	stop := make(chan struct{})
	go controller.Run(stop)
}
