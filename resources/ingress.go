package resources

import (
	"fmt"
	"github.com/gakmi/k8s-client-go/kubeconf"
	extensionsv1beat1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientExtensionsV1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/util/retry"
)

type Ingress struct {
	Name        string
	Namespace   string
	Host        string
	Path        string
	Labels      map[string]string
	ServiceName string
	ServicePort int32
}

func ingressClient(ns string) clientExtensionsV1beta1.IngressInterface {
	kubeClient, _ := kubeconf.NewKubeClient()
	//ns := "api-test"
	ingressesClient := kubeClient.Clientset.ExtensionsV1beta1().Ingresses(ns)
	return ingressesClient
}

func (s *Ingress) Create() error {
	ingressesClient := ingressClient(s.Namespace)
	// Create ingress
	fmt.Println("Creating Ingress...")
	ingress := &extensionsv1beat1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:   s.Name,
			Labels: s.Labels,
		},
		Spec: extensionsv1beat1.IngressSpec{
			Rules: []extensionsv1beat1.IngressRule{
				extensionsv1beat1.IngressRule{
					Host: s.Host,
					IngressRuleValue: extensionsv1beat1.IngressRuleValue{
						HTTP: &extensionsv1beat1.HTTPIngressRuleValue{
							Paths: []extensionsv1beat1.HTTPIngressPath{
								extensionsv1beat1.HTTPIngressPath{
									Path: s.Path,
									Backend: extensionsv1beat1.IngressBackend{
										ServiceName: s.ServiceName,
										ServicePort: intstr.IntOrString{
											IntVal: s.ServicePort,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := ingressesClient.Create(ingress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created ingress %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func (s *Ingress) Delete() error {
	ingressesClient := ingressClient(s.Namespace)
	//Delete ingress
	fmt.Println("Deleting ingress...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := ingressesClient.Delete(s.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted ingress...")
	return nil
}

func (s *Ingress) Update() error {
	ingressesClient := ingressClient(s.Namespace)
	//Update ingress
	fmt.Println("Updating ingress...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := ingressesClient.Get(s.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Ingress: %v", getErr))
		}

		result.Spec.Rules[0].Host = "register.com"
		_, updateErr := ingressesClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated ingress...")
	return nil
}

func (s *Ingress) List() error {
	ingressesClient := ingressClient(s.Namespace)
	//List ingress
	fmt.Printf("Listing ingresss in namespace %q:\n", s.Namespace)
	list, err := ingressesClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf("%s\n", d.Name)
	}
	return nil
}
