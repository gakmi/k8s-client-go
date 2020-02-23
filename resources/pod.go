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
	Node      string
	Status    string
}

//创建podsClient
func podClient(ns string) clientCoreV1.PodInterface {
	podsClient := kubeClient.Clientset.CoreV1().Pods(ns)
	return podsClient
}

//获取Pod状态
func (p *Pod) Available() bool {
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

func (p *Pod) List() []*Pod {
	podsClient := podClient(p.Namespace)
	//List pod
	fmt.Printf("Listing pods in namespace %q:\n", p.Namespace)
	list, err := podsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var pods []*Pod
	for _, d := range list.Items {
		pod := &Pod{}
		pod.Name = d.Name
		pod.Node = d.Status.HostIP
		pod.Status = string(d.Status.Phase)
		pods = append(pods, pod)
		//fmt.Printf("%s\n%s\n%s\n", d.Name, d.Status.HostIP, d.Status.Phase)
	}
	return pods
}
