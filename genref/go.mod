module github.com/kubernetes-sigs/reference-docs/genref

go 1.15

require (
	github.com/go-logr/logr v0.3.0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/tengqm/kubeconfig v0.0.0 // indirect
	github.com/yuin/goldmark v1.1.27
	github.com/yuin/goldmark-highlighting v0.0.0-20200307114337-60d527fdb691
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/tools v0.0.0-20200616133436-c1934b75d054 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/gengo v0.0.0-20201113003025-83324d819ded
	k8s.io/klog/v2 v2.4.0 // indirect
	k8s.io/kube-controller-manager v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/kube-proxy v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/kube-scheduler v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/kubelet v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/metrics v0.0.0-00010101000000-000000000000 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/tengqm/kubeconfig => github.com/tengqm/kubeconfig v0.0.0-20201104092945-d8f9a88155ff
	k8s.io/apiserver => k8s.io/apiserver v0.20.0
	k8s.io/client-go => k8s.io/client-go v0.20.0
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.20.0
	k8s.io/controller-manager => k8s.io/controller-manager v0.20.0
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.20.0
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.20.0
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.20.0
	k8s.io/kubelet => k8s.io/kubelet v0.20.0
	k8s.io/metrics => k8s.io/metrics v0.20.0
)
