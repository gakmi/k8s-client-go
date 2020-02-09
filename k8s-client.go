package main

import (
	"fmt"
	"github.com/gakmi/k8s-client-go/resources"
)

func main() {
	//	s1 := &resources.Service{
	//		Namespace:  "api-test",
	//		Name:       "register2",
	//		Port:       13001,
	//		TargetPort: 13001,
	//		NodePort:   32001,
	//		Label: map[string]string{
	//			"app": "register1",
	//		},
	//	}
	//s2 := resources.Service{
	//	Namespace: "api-test",
	//	Name:      "register2",
	//}
	//s3 := resources.Service{
	//	Namespace: "api-test",
	//}
	//s4 := resources.Service{
	//	Name:      "register2",
	//	Namespace: "api-test",
	//}
	//svc := &resources.Service{
	//	Name:      "register6",
	//	Namespace: "api-test",
	//}
	//s1.Create()
	//s2.Delete()
	//s3.List()
	//s1.GoToYaml()
	//s1.YamlToGo()
	//s4.Update(32001)

	// deployment example
	d1 := resources.Deployment{
		Namespace: "admap-test",
		Name:      "daa-imgpro22",
	}
	//d1.Create()
	//d1.List()
	//d1.Update()
	//d1.Delete()
	fmt.Println(d1.Status())

	// ingress example
	//i1 := &resources.Ingress{
	//	Namespace: "api-test",
	//	Name:      "register2",
	//}
	//i1.Create()
	//i1.Update()
	//i1.List()
	//i1.Delete()

	// configmap example
	//c1 := &resources.ConfigMap{
	//	Namespace: "api-test",
	//	Name:      "register2",
	//}
	//c1.Create()
	//c1.Update()
	//c1.List()
	//c1.Delete()
}
