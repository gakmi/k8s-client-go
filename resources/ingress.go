package resources

import (
	"fmt"
	extensionsv1beat1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientExtensionsV1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/util/retry"
	"log"
)

type Ingress struct {
	Name        string
	Namespace   string
	ServicePort int32
}

func ingressClient(ns string) clientExtensionsV1beta1.IngressInterface {
	ingressesClient := KubeClient.Clientset.ExtensionsV1beta1().Ingresses(ns)

	return ingressesClient
}

/*
创建Ingress
参数：
Name
Namespace
ServicePort
*/
func (s *Ingress) Create() error {

	var (
		ingressesClient = ingressClient(s.Namespace)
		ingress         = extensionsv1beat1.Ingress{}

		name        = s.Name
		namespace   = s.Namespace
		host        = s.Name + "." + s.Namespace
		serviceName = s.Name
		servicePort = intstr.IntOrString{
			IntVal: s.ServicePort,
		}
		backend     = &extensionsv1beat1.IngressBackend{}
		ingressRule = extensionsv1beat1.IngressRule{}
		rules       = []extensionsv1beat1.IngressRule{}
	)

	backend.ServiceName = serviceName
	backend.ServicePort = servicePort
	ingressRule.Host = host
	rules = append(rules, ingressRule)
	ingress.Name = name
	ingress.Namespace = namespace
	ingress.Spec.Rules = rules
	ingress.Spec.Backend = backend

	log.Println("Creating Ingress...")
	if result, err := ingressesClient.Create(&ingress); err != nil {
		log.Println(err)
		return err
	} else {
		log.Printf("Created ingress %q.\n", result.GetObjectMeta().GetName())
		return nil
	}

}

/*
删除Ingress
参数:
Name
Namespace
*/
func (s *Ingress) Delete() error {
	var (
		ingressesClient = ingressClient(s.Namespace)

		deletePolicy = metav1.DeletePropagationForeground
	)

	log.Println("Deleting ingress...")
	if err := ingressesClient.Delete(s.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	} else {
		log.Println("Deleted ingress...")
		return nil
	}
}

/*
更新Ingress
参数：
Name
Namespace
ServicePort
*/
func (s *Ingress) Update() error {
	var (
		ingressesClient = ingressClient(s.Namespace)
	)

	log.Println("Updating ingress...")
	if retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := ingressesClient.Get(s.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Ingress: %v", getErr))
		}

		result.Spec.Rules[0].Host = "register.com"
		_, updateErr := ingressesClient.Update(result)
		return updateErr
	}); retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	} else {
		log.Println("Updated ingress...")
		return nil
	}
}

/*
查看Ingress
参数:
Namespace
*/
func (s *Ingress) List() error {
	var (
		ingressesClient = ingressClient(s.Namespace)
	)

	log.Printf("Listing ingresss in namespace %q:\n", s.Namespace)
	if list, err := ingressesClient.List(metav1.ListOptions{}); err != nil {
		panic(err)
	} else {
		for _, d := range list.Items {
			log.Printf("%s\n", d.Name)
		}
		return nil
	}
}
