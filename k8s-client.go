package main

import (
	"github.com/gakmi/k8s-client-go/resources"
	//"time"
)

func main() {

	/*service*/
	//	var (
	//		s = resources.Service{}
	//
	//		name            = "t2"
	//		namespace       = "api-test"
	//		port      int32 = 13012
	//	)
	//
	//	s.Name = name
	//	s.Namespace = namespace
	//	s.Port = port

	//s.Delete()
	//s.Create()
	//s.List()
	//s.Update()

	/*configmap*/
	var (
		c = resources.ConfigMap{}

		name      = "t3"
		namespace = "api-test"
		logName   = "t3"
		es        = `["http://10.43.75.138:9200", "http://10.43.75.139:9200", "http://10.43.75.140:9200"]`
	)

	c.Name = name
	c.Namespace = namespace
	c.LogName = logName
	c.ES = es

	//c.List()
	//c.Update()
	c.Create()
	//c.Delete()

	/*ingress*/
	//	var (
	//		i = resources.Ingress{}
	//
	//		name              = "t2"
	//		namespace         = "api-test"
	//	)
	//
	//	i.Name = name
	//	i.Namespace = namespace

	//i.List()
	//i.Update()
	//i.Delete()
	//i.Create()

	//time.Sleep(time.Duration(6) * time.Second)
	/*deployment*/
	//	var (
	//		d = resources.Deployment{}
	//
	//		name               = "t2"
	//		namespace          = "api-test"
	//		port         int32 = 13002
	//		replicas     int32 = 2
	//		image              = "10.43.75.137/admap-dev/imp-register:20200210150618"
	//		hostNetwork        = false
	//		nodeSelector       = "node1"
	//		logDir             = "/data/logs/imp-register/"
	//		nfsDir             = "/ADMAP/IMP_DATA"
	//		nfsServer          = "10.43.75.108"
	//		nfsPath            = "/data1/admap/IMP_DATA"
	//	)
	//
	//	d.Name = name
	//	d.Namespace = namespace
	//	d.Replicas = replicas
	//	d.Image = image
	//	d.HostNetwork = hostNetwork
	//	d.Port = port
	//	d.NodeSelector = nodeSelector
	//	d.LogDir = logDir
	//	d.NFSDir = nfsDir
	//	d.NFSServer = nfsServer
	//	d.NFSPath = nfsPath
	//
	//	d.Create()
	//d.List()
	//d.Update()
	//d.Delete()

}
