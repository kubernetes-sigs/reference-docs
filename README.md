# Kubernetes Reference Docs

Tools to build reference documentation for Kubernetes APIs and CLIs.

# Api Docs

## Generate new api docs

1. From the kubernetes/kubernetes repo, copy the file `k8s.io/kubernetes/api/openapi-spec/swagger.json` to `gen_open_api/openapi-spec/swagger.json` in the reference-docs repo.

2. Update the file `gen_open_api/config.yaml`, adding any new resource definitions or operation types not already present
  - TODO: Write more on this

3. Run `make api` to build the doc html and javascript

4. Html files will be written to `gen_open_api/build`.  Copy these to where they will be hosted.

# Cli

## Generate new kubectl docs

1. Update `gen_kubectl/kubectl.yaml` by running `k8s.io/kubernetes/cmd/genkubedocs/gen_kube_docs.go` from the kubernetes/kuberentes repo and copying the file

2. Run `make cli`

3. Html files will be written to `gen_kubectl/build`.  Copy these to where they will be hosted.

# Updating brodocs version

*May need to change the image repo to one you have write access to.*

1. Update Dockerfile so it will re-clone the repo

2. Run `make brodocs`
