package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	res "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
	clientAppsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

type Deployment struct {
	Namespace    string
	Name         string
	Replicas     int32
	Image        string
	Port         int32
	HostNetwork  bool
	NodeSelector string
	NFSDir       string
	NFSServer    string
	NFSPath      string
	LogDir       string
}

func deploymentClient(ns string) clientAppsV1.DeploymentInterface {
	deploymentsClient := KubeClient.Clientset.AppsV1().Deployments(ns)

	return deploymentsClient
}

/*
创建Deployment
参数:
Name
Namespace
Replicas
Image
Port
HostNetwork
NodeSelector
NFSDir
NFSServer
NFSPath
LogDir
*/
func (s *Deployment) Create() error {

	var (
		deploymentsClient = deploymentClient(s.Namespace)
		deployment        = appsv1.Deployment{}

		name                               = s.Name
		namespace                          = s.Namespace
		replicas                           = int32Ptr(s.Replicas)
		labels                             = make(map[string]string)
		selector                           = &metav1.LabelSelector{}
		template                           = corev1.PodTemplateSpec{}
		hostNetwork                        = s.HostNetwork
		nodeSelector                       = make(map[string]string)
		containers                         = []corev1.Container{}
		containerApp                       = corev1.Container{}
		containerSidecar                   = corev1.Container{}
		port                               = corev1.ContainerPort{}
		ports                              = []corev1.ContainerPort{}
		volMLog                            = corev1.VolumeMount{}
		volMNFS                            = corev1.VolumeMount{}
		volumeMounts                       = []corev1.VolumeMount{}
		pollPolicy       corev1.PullPolicy = "IfNotPresent"
		args                               = []string{}
		valueFrom1                         = &corev1.EnvVarSource{}
		valueFrom2                         = &corev1.EnvVarSource{}
		envVar1                            = corev1.EnvVar{}
		envVar2                            = corev1.EnvVar{}
		env                                = []corev1.EnvVar{}
		runAsUser        *int64
		resourceList1    = make(map[corev1.ResourceName]res.Quantity)
		resourceList2    = make(map[corev1.ResourceName]res.Quantity)
		resources        = corev1.ResourceRequirements{}
		volMLog2         = corev1.VolumeMount{}
		volMConfig2      = corev1.VolumeMount{}
		volumeMounts2    = []corev1.VolumeMount{}
		volLog           = corev1.Volume{}
		volConfig        = corev1.Volume{}
		volNFS           = corev1.Volume{}
		volumes          = []corev1.Volume{}
		emptyDir         = &corev1.EmptyDirVolumeSource{}
		configMap        = &corev1.ConfigMapVolumeSource{}
		nfs              = &corev1.NFSVolumeSource{}
	)

	labels["app"] = s.Name
	selector.MatchLabels = labels
	if s.NodeSelector != "" {
		nodeSelector["node"] = s.NodeSelector
	}
	//containers
	port.ContainerPort = s.Port
	ports = append(ports, port)
	volMLog.Name = s.Name
	volMLog.MountPath = s.LogDir
	volMNFS.Name = s.Name + "-nfs"
	volMNFS.MountPath = s.NFSDir
	volumeMounts = append(volumeMounts, volMLog)
	volumeMounts = append(volumeMounts, volMNFS)
	containerApp.Name = s.Name
	containerApp.Image = s.Image
	containerApp.Ports = ports
	containerApp.VolumeMounts = volumeMounts
	containers = append(containers, containerApp)
	//Sidecar
	containerSidecar.Name = "filebeat"
	containerSidecar.Image = "elastic/filebeat:7.4.2"
	containerSidecar.ImagePullPolicy = pollPolicy
	args = append(args, "-c")
	args = append(args, "/etc/filebeat/filebeat.yml")
	args = append(args, "-e")
	containerSidecar.Args = args
	envVar1.Name = "POD_IP"
	valueFrom1.FieldRef = &corev1.ObjectFieldSelector{
		APIVersion: "v1",
		FieldPath:  "status.podIP",
	}
	valueFrom2.FieldRef = &corev1.ObjectFieldSelector{
		APIVersion: "v1",
		FieldPath:  "metadata.name",
	}
	envVar1.ValueFrom = valueFrom1
	envVar2.Name = "pod_name"
	envVar2.ValueFrom = valueFrom2
	env = append(env, envVar1)
	env = append(env, envVar2)
	containerSidecar.Env = env
	containerSidecar.SecurityContext = &corev1.SecurityContext{
		RunAsUser: runAsUser,
	}
	resourceList1[corev1.ResourceName("memory")] = res.MustParse("200Mi")
	resources.Limits = resourceList1
	resourceList2[corev1.ResourceName("memory")] = res.MustParse("200Mi")
	resourceList2[corev1.ResourceName("cpu")] = res.MustParse("200m")
	resources.Requests = resourceList2
	containerSidecar.Resources = resources
	volMLog2.Name = s.Name
	volMLog2.MountPath = "/log"
	volMConfig2.Name = "config"
	volMConfig2.MountPath = "/etc/filebeat/"
	volumeMounts2 = append(volumeMounts2, volMLog2)
	volumeMounts2 = append(volumeMounts2, volMConfig2)
	containerSidecar.VolumeMounts = volumeMounts2
	containers = append(containers, containerSidecar)

	volLog.Name = s.Name
	volLog.VolumeSource.EmptyDir = emptyDir
	volConfig.Name = "config"
	configMap.Name = s.Name
	volNFS.Name = s.Name + "-nfs"
	nfs.Server = s.NFSServer
	nfs.Path = s.NFSPath
	volNFS.VolumeSource.NFS = nfs
	volConfig.ConfigMap = configMap
	volumes = append(volumes, volLog)
	volumes = append(volumes, volConfig)
	volumes = append(volumes, volNFS)
	//template
	template.Labels = labels
	template.Spec.HostNetwork = hostNetwork
	if nodeSelector != nil {
		template.Spec.NodeSelector = nodeSelector
	}
	template.Spec.Containers = containers
	template.Spec.Volumes = volumes

	deployment.Name = name
	deployment.Namespace = namespace
	deployment.Labels = labels

	deployment.Spec.Replicas = replicas
	deployment.Spec.Selector = selector
	deployment.Spec.Template = template

	log.Println("Creating Deployment...")
	if result, err := deploymentsClient.Create(&deployment); err != nil {
		log.Println(err)
		return err
	} else {
		log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
		return nil
	}

}

/*
删除Deployment
参数:
Name
Namespace
*/
func (s *Deployment) Delete() error {
	var (
		deploymentsClient = deploymentClient(s.Namespace)
	)

	deletePolicy := metav1.DeletePropagationForeground
	log.Println("Deleting deployment...")
	if err := deploymentsClient.Delete(s.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	} else {
		log.Println("Deleted deployment.")
		return nil
	}
}

/*
更新Deployment
参数:
Name
Namespace
*/
func (s *Deployment) Update() error {

	var (
		deploymentsClient = deploymentClient(s.Namespace)
	)

	log.Println("Updating deployment...")
	if retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(s.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13"
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	}); retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	} else {
		log.Println("Updated deployment...")
		return nil
	}
}

/*
查看Deployment
参数:
Namespace
*/
func (s *Deployment) List() error {
	var (
		deploymentsClient = deploymentClient(s.Namespace)
	)

	log.Printf("Listing Deployment in namespace %q:\n", s.Namespace)
	if list, err := deploymentsClient.List(metav1.ListOptions{}); err != nil {
		panic(err)
	} else {
		for _, d := range list.Items {
			log.Printf("%s\n", d.Name)
		}
		return nil
	}
}

/*
获取Pod状态
*/
func (s *Deployment) Status() bool {
	var (
		deploymentsClient = deploymentClient(s.Namespace)
	)
	//status deployment
	log.Println("status deployment...")
	if deployment, err := deploymentsClient.Get(s.Name, metav1.GetOptions{}); err != nil {
		log.Println(err)
		return false
	} else {
		log.Println(deployment.Status.UnavailableReplicas)
		if deployment.Status.UnavailableReplicas > 0 {
			return false
		}
		return true
	}
}

func int32Ptr(i int32) *int32 { return &i }

func ContainerSidecar() appsv1.Deployment {

	var (
		containerYaml []byte
		containerJson []byte
		//container     = corev1.Container{}
		deployment = appsv1.Deployment{}
		err        error
	)
	// 读取YAML
	if containerYaml, err = ioutil.ReadFile("./yamls/deployment.yaml"); err != nil {
		log.Println(err)
		log.Println(containerYaml)
	}
	// YAML转JSON
	if containerJson, err = yaml2.ToJSON(containerYaml); err != nil {
		log.Println(err)
		log.Println(containerJson)
	}

	// JSON转struct
	if err = json.Unmarshal(containerJson, &deployment); err != nil {
		log.Println(err)
	}

	log.Println(deployment)
	return deployment

}
