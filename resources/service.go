package resources

import (
	"fmt"
	"github/gakmi/k8s-client-go/kubeconf"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

type Service struct {
	Namespace  string
	Name       string
	Port       int32
	NodePort   int32
	TargetPort int32
	Label      map[string]string
}

//创建servicesClient
func client(ns string) clientCoreV1.ServiceInterface {
	kubeClient, _ := kubeconf.NewKubeClient()
	//ns := "api-test"
	servicesClient := kubeClient.Clientset.CoreV1().Services(ns)
	return servicesClient
}

func (s *Service) Create() error {
	servicesClient := client(s.Namespace)
	// Create service
	fmt.Println("Creating Service...")
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Name,
			Labels: map[string]string{
				"project": "api-test",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     s.Name,
					Protocol: corev1.ProtocolTCP,
					Port:     s.Port,
					TargetPort: intstr.IntOrString{
						IntVal: s.TargetPort,
					},
					NodePort: s.NodePort,
				},
			},
			Selector: s.Label,
			Type:     corev1.ServiceTypeNodePort,
		},
	}
	result, err := servicesClient.Create(service)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func (s *Service) Delete() error {
	servicesClient := client(s.Namespace)
	//Delete service
	fmt.Println("Deleting service...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := servicesClient.Delete(s.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
	return nil
}

func (s *Service) Update(nodeport int32) error {
	servicesClient := client(s.Namespace)
	//Update service
	fmt.Println("Updating service...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := servicesClient.Get(s.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Ports[0].NodePort = nodeport
		_, updateErr := servicesClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated service...")
	return nil
}

func (s *Service) List() error {
	servicesClient := client(s.Namespace)
	//List service
	fmt.Printf("Listing services in namespace %q:\n", s.Namespace)
	list, err := servicesClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf("%s\n", d.Name)
	}
	return nil
}
