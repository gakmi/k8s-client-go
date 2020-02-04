package kubeconf

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

//kubeClient结构
type KubeClient struct {
	kubeClientConfig *restclient.Config
	Clientset        *kubernetes.Clientset
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

/*
创建kube config
*/
func newKubeConfig() *restclient.Config {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	return config
}

/*
创建clientset
*/
func (kc *KubeClient) Create() error {
	if kc.Clientset == nil {
		clientset, err := kubernetes.NewForConfig(kc.kubeClientConfig)
		if err != nil {
			return err
		}
		kc.Clientset = clientset
	}

	return nil
}

/*
创建KubeClient
*/
func NewKubeClient() (*KubeClient, error) {
	kubeClient := &KubeClient{kubeClientConfig: newKubeConfig()}
	if err := kubeClient.Create(); err != nil {
		return nil, err
	}

	return kubeClient, nil
}
