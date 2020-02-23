package resources

import (
	"github.com/gakmi/k8s-client-go/kubeconf"
)

var KubeClient *kubeconf.KubeClient

func init() {
	var err error
	if KubeClient, err = kubeconf.NewKubeClient(); err != nil {
		panic(err)
	}
}
