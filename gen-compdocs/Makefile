all:
	mkdir -p build
	go run main.go build kube-apiserver
	go run main.go build kube-controller-manager
	go run main.go build kube-scheduler
	go run proxy/main.go build
	go run main.go build kubeadm
	go run main.go build kubectl
