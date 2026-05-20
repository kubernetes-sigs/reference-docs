# gen-compdocs

`gen-compdocs` generates Markdown reference pages for the Kubernetes command-line tools. The generator imports the cobra command trees from `k8s.io/kubernetes/cmd/*` and `k8s.io/kubectl/pkg/cmd`, then renders each command and subcommand as a Hugo-compatible Markdown page.

Modules covered:

- `kube-apiserver`
- `kube-controller-manager`
- `kube-scheduler`
- `kube-proxy`
- `kubelet`
- `kubeadm` (and subcommands)
- `kubectl` (and subcommands)

`gen-compdocs` is the replacement for `gen-kubectldocs`.

For the canonical end-to-end walkthrough, see [Generating Reference Documentation for kubectl Commands](https://kubernetes.io/docs/contribute/generate-ref-docs/kubectl/).

## Inputs

| Input | Purpose |
|---|---|
| `K8S_RELEASE` | Kubernetes release version, for example `1.36.0`. |
| `K8S_WEBROOT` | Path to a local `kubernetes/website` checkout. Required for the copy targets. |
| `gen-compdocs/go.mod` | Pins the Kubernetes source versions that supply the command definitions. Update with `go get` and `go mod tidy` for a new release. |

## Build

From the repository root:

```shell
make comp
```

This delegates to `gen-compdocs/Makefile`, which runs:

```shell
go run main.go build <module>
```

for each module listed above, writing output to `gen-compdocs/build/`.

## Copy to website

Set `K8S_WEBROOT` to a `kubernetes/website` checkout, then run one or more of:

```shell
make copycomp-core      # kube-apiserver, kube-controller-manager, kube-scheduler, kube-proxy, kubelet
make copycomp-kubeadm   # kubeadm and subcommands
make copycomp-kubectl   # kubectl and subcommands
make copycomp           # all of the above
```

Each `copycomp-*` target overlays files onto the destination directory. The destination keeps hand-curated `_index.md` (and a `README.md` under `kubeadm/generated/`); these files are preserved across regenerations.

## Output

- `gen-compdocs/build/<module>.md` — top-level command page.
- `gen-compdocs/build/<module>_<subcommand>.md` — subcommand pages.

