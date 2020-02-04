package main

import (
	//"fmt"
	"k8s-client-go/resources"
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
	s4 := resources.Service{
		Name:      "register2",
		Namespace: "api-test",
	}
	//svc := &resources.Service{
	//	Name:      "register6",
	//	Namespace: "api-test",
	//}
	//s1.Create()
	//s2.Delete()
	//s3.List()
	s4.Update(32001)
}
