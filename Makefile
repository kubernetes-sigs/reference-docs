K8SIOROOT=../../../../../go/src/k8s.io/kubernetes.github.io
K8SROOT=../../../../../go/src/k8s.io/kubernetes

all: main.go
	go build main.go

cleancli:
	rm -f main
	rm -rf $(shell pwd)/gen_kubectl/includes
	rm -rf $(shell pwd)/gen_kubectl/build
	rm -rf $(shell pwd)/gen_kubectl/manifest.json

cleanapi:
	rm -f main
	rm -rf $(shell pwd)/gen_open_api/build
	rm -rf $(shell pwd)/gen_open_api/includes
	rm -rf $(shell pwd)/gen_open_api/manifest.json
	rm -rf $(shell pwd)/gen_open_api/includes/_generated_*

brodocs:
	docker build . -t pwittrock/brodocs
	docker push pwittrock/brodocs

cli: cleancli
	go run main.go --doc-type kubectl --kubernetes-version v1_6
	docker run -v $(shell pwd)/gen_kubectl/includes:/source -v $(shell pwd)/gen_kubectl/build:/build -v $(shell pwd)/gen_kubectl/:/manifest pwittrock/brodocs

// Usage: TAG="vN" make pushcli
pushcli: cli
	cd gen_kubectl && docker build . -t pwittrock/cli-docs:$(TAG)
	docker push pwittrock/cli-docs:$(TAG)

copycli: cli
	rm -rf gen_kubectl/build/documents/
	rm -rf gen_kubectl/build/runbrodocs.sh
	rm -rf gen_kubectl/build/manifest.json
	rm -rf $(K8SIOROOT)/docs/user-guide/kubectl/v1.6/*
	cp -r gen_kubectl/build/* $(K8SIOROOT)/docs/user-guide/kubectl/v1.6/

pushcliconfig:
	cd configs && kubectl apply -f kubectldocs.yaml

api: cleanapi
	go run main.go --doc-type open-api
	docker run -v $(shell pwd)/gen_open_api/includes:/source -v $(shell pwd)/gen_open_api/build:/build -v $(shell pwd)/gen_open_api/:/manifest pwittrock/brodocs

// Usage: TAG="vN" make pushapi
pushapi: api
	cd gen_open_api && docker build . -t pwittrock/api-docs:$(TAG)
	docker push pwittrock/api-docs:$(TAG)

updateapispec: api
	cp $(K8SROOT)/api/openapi-spec/swagger.json gen_open_api/openapi-spec/swagger.json

ca:
	ls $(K8SIOROOT) /docs/api-reference/v1.6/

copyapi: api
	rm -rf gen_open_api/build/documents/
	rm -rf gen_open_api/build/runbrodocs.sh
	rm -rf gen_open_api/build/manifest.json
	rm -rf $(K8SIOROOT)/docs/api-reference/v1.6/*
	cp -r gen_open_api/build/* $(K8SIOROOT)/docs/api-reference/v1.6/

pushapiconfig:
	cd configs && kubectl apply -f apidocs.yaml

resource: cleanapi
	go run main.go --doc-type open-api  --build-operations=false
	docker run -v $(shell pwd)/gen_open_api/includes:/source -v $(shell pwd)/gen_open_api/build:/build -v $(shell pwd)/gen_open_api/:/manifest pwittrock/brodocs

// Usage: TAG="vN" make pushresource
pushresource: resource
	cd gen_open_api && docker build . -t pwittrock/resource-docs:$(TAG)
	docker push pwittrock/resource-docs:$(TAG)

copyresource: resource
	rm -rf gen_open_api/build/documents/
	rm -rf gen_open_api/build/runbrodocs.sh
	rm -rf gen_open_api/build/manifest.json
	rm -rf $(K8SIOROOT)/docs/resources-reference/v1.6/*
	cp -r gen_open_api/build/* $(K8SIOROOT)/docs/resources-reference/v1.6/

pushresourceconfig:
	cd configs && kubectl apply -f resourcedocs.yaml

pushstagingapiconfig:
	cd configs && kubectl apply -f stagingapidocs.yaml
