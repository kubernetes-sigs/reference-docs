package main

// synthetic dependencies to prevent `go mod tidy`
// from "accidentally" removing them
import (
	_ "k8s.io/api"
	_ "k8s.io/apimachinery"
	_ "k8s.io/apiserver"
	_ "k8s.io/client-go"
	_ "k8s.io/cloud-provider"
	_ "k8s.io/cluster-bootstrap"
	_ "k8s.io/component-base"
	_ "k8s.io/controller-manager"
	_ "k8s.io/kube-controller-manager"
	_ "k8s.io/kube-proxy"
	_ "k8s.io/kube-scheduler"
	_ "k8s.io/kubelet"
	_ "k8s.io/metrics"
)
