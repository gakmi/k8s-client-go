package resources

import (
	"fmt"
	"github.com/gakmi/k8s-client-go/kubeconf"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

type ConfigMap struct {
	Namespace string
	Name      string
	LogName   string
}

func configmapClient(ns string) clientCoreV1.ConfigMapInterface {
	kubeClient, _ := kubeconf.NewKubeClient()
	configmapsClient := kubeClient.Clientset.CoreV1().ConfigMaps(ns)

	return configmapsClient
}

/*
创建ConfigMap
参数：
Name
Namespace
LogName
*/
func (s *ConfigMap) Create() error {

	var (
		configmapsClient = configmapClient(s.Namespace)
		configmap        = corev1.ConfigMap{}

		name      = s.Name
		namespace = s.Namespace
		labels    = make(map[string]string)
		data      = make(map[string]string)
	)

	labels["app"] = s.Name
	data["filebeat.yml"] = fmt.Sprintf(`filebeat.inputs:
- type: log
  paths:
    - "/log/%s-${POD_IP}.json"
  tail_files: true
  fields:
    pod_name: '${pod_name}'
    POD_IP: '${POD_IP}'
  tags: ["%s"]
  json.keys_under_root: true
  json.overwrite_keys: true
  json.message_key: message
  exclude_files: ['\.gz$']
setup.template.name: "%s-%s"
setup.template.pattern: "%s-%s-*"
setup.ilm.enabled: false
output.elasticsearch:
  hosts: ["http://10.43.75.138:9200", "http://10.43.75.139:9200", "http://10.43.75.140:9200"]
  index: "%s-%s-%s"
processors:
  - drop_fields:
      fields: ["input","log","beat","offset","source","host","span","trace","parent"]`, s.LogName, s.Name, s.Name, s.Namespace, s.Name, s.Namespace, s.Name, s.Namespace, "%{+yyyy.MM.dd}")

	configmap.Name = name
	configmap.Namespace = namespace
	configmap.Labels = labels
	configmap.Data = data

	log.Println("Creating ConfigMap...")
	if result, err := configmapsClient.Create(&configmap); err != nil {
		log.Println(err)
		return err
	} else {
		log.Printf("Created configmap %q.\n", result.GetObjectMeta().GetName())
		return nil
	}

}

/*
删除ConfigMap
参数：
Name
Namespace
*/
func (s *ConfigMap) Delete() error {
	var (
		configmapsClient = configmapClient(s.Namespace)
		deletePolicy     = metav1.DeletePropagationForeground
	)

	log.Println("Deleting configmap...")
	if err := configmapsClient.Delete(s.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	} else {
		log.Println("Deleted configmap.")
		return nil
	}
}

/*
更新ConfigMap
参数：
Name
Namespace
LogName
*/
func (s *ConfigMap) Update() error {
	var (
		configmapsClient = configmapClient(s.Namespace)
	)
	log.Println("Updating configmap...")
	if retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := configmapsClient.Get(s.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Data["filebeat.yml"] = fmt.Sprintf(`filebeat.inputs:
- type: log
  paths:
    - "/log/%s-${POD_IP}.json"
  tail_files: true
  fields:
    pod_name: '${pod_name}'
    POD_IP: '${POD_IP}'
  tags: ["%s"]
  json.keys_under_root: true
  json.overwrite_keys: true
  json.message_key: message
  exclude_files: ['\.gz$']
setup.template.name: "%s-%s"
setup.template.pattern: "%s-%s-*"
setup.ilm.enabled: false
output.elasticsearch:
  hosts: ["http://10.43.75.138:9200", "http://10.43.75.139:9200", "http://10.43.75.140:9200"]
  index: "%s-%s-%s"
processors:
  - drop_fields:
      fields: ["input","log","beat","offset","source","host","span","trace","parent"]`, s.Name, s.Name, s.Name, s.Namespace, s.Name, s.Namespace, s.Name, s.Namespace, "%{+yyyy.MM.dd}")

		_, updateErr := configmapsClient.Update(result)
		return updateErr
	}); retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	} else {
		log.Println("Updated configmap...")
		return nil
	}
}

/*
查看ConfigMap
参数：
Namespace
*/
func (s *ConfigMap) List() error {
	var (
		configmapsClient = configmapClient(s.Namespace)
	)

	log.Printf("Listing ConfigMap in namespace %q:\n", s.Namespace)
	if list, err := configmapsClient.List(metav1.ListOptions{}); err != nil {
		panic(err)
	} else {
		for _, d := range list.Items {
			log.Printf("%s\n", d.Name)
		}
		return nil
	}
}
