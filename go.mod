module k8s-client-go

go 1.13

require (
	k8s.io/api v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v0.16.4
)

replace (
	k8s-client-go/kubeconf => /work/pj/k8s-client-go/kubeconf
	k8s-client-go/resources => /work/pj/k8s-client-go/resources
)
