# To generate docs with make targets
# from this repo base directory,
# set the following environment variables
# to match your environment and release version
#
# WEBROOT=~/src/github.com/kubernetes/website
# K8SROOT=~/k8s/src/k8s.io/kubernetes
# K8S_RELEASE=1.17

WEBROOT=${K8S_WEBROOT}
K8SROOT=${K8S_ROOT}
K8SRELEASE=${K8S_RELEASE}

# create a directory name from release string: 1.17 -> 1_17
K8SRELEASEDIR=$(shell echo "$(K8SRELEASE)" | sed "s/\./_/g")

APISRC=gen-apidocs
APIDST=$(WEBROOT)/static/docs/reference/generated/kubernetes-api/v$(K8SRELEASE)

CLISRC=gen-kubectldocs
CLIDST=$(WEBROOT)/static/docs/reference/generated/kubectl

default:
	@echo "Support commands:\napi cli comp copyapi copycli createversiondirs updateapispec"

# create directories for new release
createversiondirs:
	@echo "Calling set_version_dirs.sh"
	./set_version_dirs.sh
	@echo "K8S Release dir: $(K8SRELEASEDIR)"

# Build kubectl docs
cleancli:
	rm -rf $(shell pwd)/gen-kubectldocs/build

cli: cleancli
	go run gen-kubectldocs/main.go --kubernetes-release=$(K8SRELEASE)

copycli: cli
# make a versioned directory?
	cp $(CLISRC)/build/kubectl-commands.html $(CLIDST)/kubectl-commands.html

	# copy js files
	mkdir -p $(CLIDST)/js
	cp $(CLISRC)/build/navData.js $(CLIDST)/js/navData.js
	cp $(CLISRC)/static/js/* $(CLIDST)/js/

	# copy css files
	mkdir -p $(CLIDST)/css
	cp $(CLISRC)/static/css/* $(CLIDST)/css/

	# copy fonts data
	mkdir -p $(CLIDST)/fonts
	cp $(CLISRC)/static/fonts/* $(CLIDST)/fonts/

# Build kube component,tool docs
cleancomp:
	rm -rf $(shell pwd)/gen-compdocs/build

comp: cleancomp
	mkdir -p gen-compdocs/build
	go run gen-compdocs/main.go gen-compdocs/build kube-apiserver
	go run gen-compdocs/main.go gen-compdocs/build kube-controller-manager
	go run gen-compdocs/main.go gen-compdocs/build cloud-controller-manager
	go run gen-compdocs/main.go gen-compdocs/build kube-scheduler
	go run gen-compdocs/main.go gen-compdocs/build kubelet
	go run gen-compdocs/main.go gen-compdocs/build kube-proxy
	go run gen-compdocs/main.go gen-compdocs/build kubeadm
	go run gen-compdocs/main.go gen-compdocs/build kubectl

# Build api docs
# Note: May want to clean dir every time and fetch new copy of swagger.json
# validate size > 0 of swagger.json
updateapispec: createversiondirs
	@echo "Updating swagger.json for release v$(K8SRELEASE).0"
	CURDIR=$(shell pwd)
	if ! [ -f $(APISRC)/config/v$(K8SRELEASEDIR)/swagger.json ]; then \
		cd $(K8SROOT); \
		git show "v$(K8SRELEASE).0:api/openapi-spec/swagger.json" > swagger.json.$(K8SRELEASE); \
		mv swagger.json.$(K8SRELEASE) $(CURDIR)/$(APISRC)/config/v$(K8SRELEASEDIR)/swagger.json; \
		cd $(CURDIR); \
		echo Current dir $(shell pwd); \
	fi

api: cleanapi
	go run gen-apidocs/main.go --kubernetes-release=$(K8SRELEASE) --work-dir=gen-apidocs --munge-groups=false

cleanapi:
	rm -rf $(shell pwd)/gen-apidocs/build

copyapi: api
	mkdir -p $(APIDST)
	cp $(APISRC)/build/index.html $(APIDST)/index.html
	# copy scroll.js, jquery.scrollTo.min.js and the new navData.js
	mkdir -p $(APIDST)/js
	cp $(APISRC)/build/navData.js $(APIDST)/js/
	cp $(APISRC)/static/js/* $(APIDST)/js/
	# copy stylesheet.css, bootstrap.min.css, font-awesome.min.css
	mkdir -p $(APIDST)/css
	cp $(APISRC)/static/css/* $(APIDST)/css/
	# copy fonts data
	mkdir -p $(APIDST)/fonts
	cp $(APISRC)/static/fonts/* $(APIDST)/fonts/
