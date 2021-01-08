module github.com/kubernetes-sigs/reference-docs/gen-compdocs

go 1.15

require (
	github.com/go-logr/logr v0.3.0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/yuin/goldmark v1.1.27 // indirect
	github.com/yuin/goldmark-highlighting v0.0.0-20200307114337-60d527fdb691 // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/tools v0.0.0-20200616133436-c1934b75d054 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/gengo v0.0.0-20201113003025-83324d819ded // indirect
	k8s.io/klog/v2 v2.4.0 // indirect
	k8s.io/kubectl v0.0.0
	k8s.io/kubernetes v0.0.0
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.20.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.0
	k8s.io/apiserver => k8s.io/apiserver v0.20.0
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.0
	k8s.io/client-go => k8s.io/client-go v0.20.0
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.20.0
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.0
	k8s.io/code-generator => k8s.io/code-generator v0.20.0
	k8s.io/component-base => k8s.io/component-base v0.20.0
	k8s.io/component-helpers => k8s.io/component-helpers v0.20.0
	k8s.io/controller-manager => k8s.io/controller-manager v0.20.0
	k8s.io/cri-api => k8s.io/cri-api v0.20.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.20.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.20.0
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.20.0
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.20.0
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.20.0
	k8s.io/kubectl => k8s.io/kubectl v0.20.0
	k8s.io/kubelet => k8s.io/kubelet v0.20.0
	k8s.io/kubernetes => ../../../../k8s.io/kubernetes
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.20.0
	k8s.io/metrics => k8s.io/metrics v0.20.0
	k8s.io/mount-utils => k8s.io/mount-utils v0.20.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.20.0
)
