package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//"k8s.io/client-go/util/retry"
	"log"
)

type Namespace struct{}

//创建nodesClient
func namespaceClient() clientCoreV1.NamespaceInterface {
	namespacesClient := KubeClient.Clientset.CoreV1().Namespaces()

	return namespacesClient
}

//获取Node
func (s *Namespace) List() ([]corev1.Namespace, error) {

	var (
		namespacesClient = namespaceClient()
	)

	log.Printf("Listing Namespace:\n")
	if list, err := namespacesClient.List(metav1.ListOptions{}); err != nil {
		panic(err)
	} else {
		for _, d := range list.Items {
			log.Printf("%s\n", d.Name)
		}
		return list.Items, nil
	}

}
