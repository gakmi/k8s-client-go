package resources

import (
	"fmt"
	"github.com/gakmi/k8s-client-go/kubeconf"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

type Service struct {
	Namespace string
	Name      string
	Port      int32
}

//创建servicesClient
func serviceClient(ns string) clientCoreV1.ServiceInterface {
	servicesClient := kubeClient.Clientset.CoreV1().Services(ns)

	return servicesClient

}

/*
创建Service
参数：
Name
Namespace
Port
*/
func (s *Service) Create() error {

	var (
		servicesClient = serviceClient(s.Namespace)
		service        = corev1.Service{}

		name        = s.Name
		namespace   = s.Namespace
		ports       = []corev1.ServicePort{}
		labels      = make(map[string]string)
		annotations = make(map[string]string)
		selector    = make(map[string]string)
		servicePort corev1.ServicePort
		serviceType corev1.ServiceType
	)

	servicePort.TargetPort = intstr.IntOrString{
		IntVal: s.Port,
	}
	servicePort.Port = s.Port
	servicePort.Name = s.Name
	ports = append(ports, servicePort)
	labels["app"] = s.Name
	selector = labels
	serviceType = "ClusterIP"

	//设置labels
	service.Labels = labels
	//设置annotations
	service.Annotations = annotations
	//设置name
	service.Name = name
	//设置namespace
	service.Namespace = namespace
	//设置ports
	service.Spec.Ports = ports
	//设置selector
	service.Spec.Selector = selector
	//设置type
	service.Spec.Type = serviceType

	log.Println("Createing Service...")
	if _, err := servicesClient.Create(&service); err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Created Service.")
		return nil
	}

}

/*
删除Service
参数：
Name
Namespace
*/
func (s *Service) Delete() error {

	var (
		servicesClient = serviceClient(s.Namespace)

		deletePolicy = metav1.DeletePropagationForeground
	)

	log.Println("Deleting service...")
	if err := servicesClient.Delete(s.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Deleted service.")
		return nil
	}

}

/*
更新Service
Name
Namespace
Port
*/
func (s *Service) Update() error {

	var (
		servicesClient = serviceClient(s.Namespace)
	)

	log.Println("Updating service...")
	if retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := servicesClient.Get(s.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Ports[0].Port = s.Port
		result.Spec.Ports[0].TargetPort = intstr.IntOrString{
			IntVal: s.Port,
		}
		_, updateErr := servicesClient.Update(result)

		return updateErr
	}); retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	} else {
		log.Println("Updated service...")
		return nil
	}

}

/*
查看Service
参数：
Namespace
*/
func (s *Service) List() error {

	var (
		servicesClient = serviceClient(s.Namespace)
	)

	log.Printf("Listing Service in namespace %q:\n", s.Namespace)
	if list, err := servicesClient.List(metav1.ListOptions{}); err != nil {
		panic(err)
	} else {
		for _, d := range list.Items {
			log.Printf("%s\n", d.Name)
		}
		return nil
	}

}
