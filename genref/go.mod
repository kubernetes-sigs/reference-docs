module github.com/kubernetes-sigs/reference-docs/genref

go 1.16

require (
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/tengqm/kubeconfig v0.0.0-20220708134523-a4b352bcb4fd // indirect
	github.com/yuin/goldmark v1.4.1
	github.com/yuin/goldmark-highlighting v0.0.0-20210516132338-9216f9c5aa01
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/tools v0.1.10-0.20220218145154-897bd77cd717 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.24.0 // indirect
	k8s.io/apiserver v0.24.0
	k8s.io/client-go v0.24.0
	k8s.io/cloud-provider v0.24.0
	k8s.io/cluster-bootstrap v0.24.0
	k8s.io/component-base v0.24.0
	k8s.io/controller-manager v0.24.0
	k8s.io/gengo v0.0.0-20211129171323-c02415ce4185
	k8s.io/klog/v2 v2.60.1
	k8s.io/kube-controller-manager v0.24.0
	k8s.io/kube-proxy v0.24.0
	k8s.io/kube-scheduler v0.24.0
	k8s.io/kubelet v0.24.0
	k8s.io/metrics v0.24.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f
)
