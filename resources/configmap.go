package resources

import (
	"fmt"
	"github.com/gakmi/k8s-client-go/kubeconf"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

type ConfigMap struct {
	Namespace string
	Name      string
	Labels    map[string]string
	Data      map[string]string
}

func configmapClient(ns string) clientCoreV1.ConfigMapInterface {
	kubeClient, _ := kubeconf.NewKubeClient()
	//ns := "api-test"
	configmapsClient := kubeClient.Clientset.CoreV1().ConfigMaps(ns)
	return configmapsClient
}

func (s *ConfigMap) Create() error {
	configmapsClient := configmapClient(s.Namespace)
	// Create configmap
	fmt.Println("Creating ConfigMap...")
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   s.Name,
			Labels: s.Labels,
		},
		Data: s.Data,
	}
	result, err := configmapsClient.Create(configmap)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created configmap %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func (s *ConfigMap) Delete() error {
	configmapsClient := configmapClient(s.Namespace)
	//Delete configmap
	fmt.Println("Deleting configmap...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := configmapsClient.Delete(s.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted configmap.")
	return nil
}

func (s *ConfigMap) Update() error {
	configmapsClient := configmapClient(s.Namespace)
	//Update configmap
	fmt.Println("Updating configmap...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := configmapsClient.Get(s.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Data["filebeat.yml"] =
			`filebeat.inputs:
    - type: log
      paths:
        - "/log/register-${POD_IP}.json"
      tail_files: true
      fields:
        pod_name: '${pod_name}'
        POD_IP: '${POD_IP}'
      tags: ["register"]
      json.keys_under_root: true
      json.overwrite_keys: true
      json.message_key: message
      exclude_files: ['\.gz$']
    setup.template.name: "register-test"
    setup.template.pattern: "register-test-*"
    setup.ilm.enabled: false
    output.elasticsearch:
      hosts: ["http://10.43.75.138:9200", "http://10.43.75.139:9200", "http://10.43.75.140:9200"]
      index: "register-test-%{+yyyy.MM.dd}"
    processors:
      - drop_fields:
          fields: ["input","log","beat","offset","source","host","span","trace","parent"]`
		_, updateErr := configmapsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated configmap...")
	return nil
}

func (s *ConfigMap) List() error {
	configmapsClient := configmapClient(s.Namespace)
	//List configmap
	fmt.Printf("Listing configmaps in namespace %q:\n", s.Namespace)
	list, err := configmapsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf("%s\n", d.Name)
	}
	return nil
}
