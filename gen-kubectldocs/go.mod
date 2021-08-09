module github.com/kubernetes-sigs/reference-docs/gen-kubectldocs

go 1.16

require (
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/googleapis/gnostic v0.5.5
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/cli-runtime v0.22.0
	k8s.io/client-go v0.22.0
	k8s.io/component-base v0.22.0
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e
	k8s.io/kubectl v0.22.0
	k8s.io/metrics v0.22.0
	k8s.io/utils v0.0.0-20210707171843-4b05e18ac7d9
)

replace (
	github.com/jonboulle/clockwork => github.com/jonboulle/clockwork v0.2.2
	github.com/moby/term => github.com/moby/term v0.0.0-20210610120745-9d4ed1856297
	sigs.k8s.io/kustomize/kustomize/v4 => sigs.k8s.io/kustomize/kustomize/v4 v4.2.0
	sigs.k8s.io/kustomize/kyaml => sigs.k8s.io/kustomize/kyaml v0.11.0
)
