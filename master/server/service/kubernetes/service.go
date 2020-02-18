package kubernetes

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func GetAllServices() ([]apiv1.Service, error) {
	services, err := ClientSet.CoreV1().Services(config.KubernetesConf.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return services.Items, nil
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

func GetServiceTaskID(service *apiv1.Service) (uint64, error) {
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
