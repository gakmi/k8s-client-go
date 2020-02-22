module github.com/gakmi/k8s-client-go

go 1.13

replace (
	github.com/gakmi/k8s-client-go/kubeconf => /opt/k8s-client-go/kubeconf
	github.com/gakmi/k8s-client-go/resources => /opt/k8s-client-go/resources
)

require (
	github.com/imdario/mergo v0.3.8 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	k8s.io/api v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v0.16.4
	k8s.io/utils v0.0.0-20200124190032-861946025e34 // indirect
)
