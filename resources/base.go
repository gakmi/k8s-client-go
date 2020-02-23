package resources

import "github.com/gakmi/k8s-client-go/kubeconf"

var kubeClient kubeconf.KubeClient

func init() {
	if kubeClient, err := kubeconf.NewKubeClient(); err != nil {
		panic(err)
	}
}
