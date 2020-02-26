package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//"k8s.io/client-go/util/retry"
	"log"
)

type Node struct{}

//创建nodesClient
func nodeClient() clientCoreV1.NodeInterface {
	nodesClient := KubeClient.Clientset.CoreV1().Nodes()

	return nodesClient
}

//获取Node
func (s *Node) List() ([]corev1.Node, error) {

	var (
		nodesClient = nodeClient()
	)

	log.Printf("Listing Node:\n")
	if list, err := nodesClient.List(metav1.ListOptions{}); err != nil {
		panic(err)
	} else {
		for _, d := range list.Items {
			log.Printf("%s\n", d.Name)
		}
		return list.Items, nil
	}

}
