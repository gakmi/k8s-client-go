package resources

import (
	"fmt"
	"github.com/gakmi/k8s-client-go/kubeconf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Pod struct {
	Namespace string
	Name      string
}

//创建podsClient
func podClient(ns string) clientCoreV1.PodInterface {
	kubeClient, _ := kubeconf.NewKubeClient()
	//ns := "api-test"
	podsClient := kubeClient.Clientset.CoreV1().Pods(ns)
	return podsClient
}

//获取Pod状态
func (p *Pod) Status() bool {
	podsClient := podClient(p.Namespace)
	//status pod
	fmt.Println("status pod...")
	pod, err := podsClient.Get(p.Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
		return false
	}
	fmt.Println(pod)
	return true
}
