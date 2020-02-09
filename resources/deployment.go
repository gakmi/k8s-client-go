package resources

import (
	"fmt"
	"github.com/gakmi/k8s-client-go/kubeconf"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientAppsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
)

type Deployment struct {
	Namespace      string
	Name           string
	Replicas       int32
	Labels         map[string]string
	ContainerName  string
	ContainerImage string
	ContainerPort  int32
}

func deploymentClient(ns string) clientAppsV1.DeploymentInterface {
	kubeClient, _ := kubeconf.NewKubeClient()
	//ns := "api-test"
	deploymentsClient := kubeClient.Clientset.AppsV1().Deployments(ns)
	return deploymentsClient
}

func (s *Deployment) Create() error {
	deploymentsClient := deploymentClient(s.Namespace)
	// Create deployment
	fmt.Println("Creating Deployment...")
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(s.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: s.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: s.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  s.ContainerName,
							Image: s.ContainerImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: s.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}

func (s *Deployment) Delete() error {
	deploymentsClient := deploymentClient(s.Namespace)
	//Delete deployment
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(s.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
	return nil
}

func (s *Deployment) Update() error {
	deploymentsClient := deploymentClient(s.Namespace)
	//Update deployment
	fmt.Println("Updating deployment...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(s.Name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13"
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")
	return nil
}

func (s *Deployment) List() error {
	deploymentsClient := deploymentClient(s.Namespace)
	//List deployment
	fmt.Printf("Listing deployments in namespace %q:\n", s.Namespace)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf("%s\n", d.Name)
	}
	return nil
}

//获取Pod状态
func (s *Deployment) Status() bool {
	deploymentsClient := deploymentClient(s.Namespace)
	//status deployment
	fmt.Println("status deployment...")
	deployment, err := deploymentsClient.Get(s.Name, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(deployment.Status.UnavailableReplicas)
	if deployment.Status.UnavailableReplicas > 0 {
		return false
	}
	return true
}

func int32Ptr(i int32) *int32 { return &i }
