# Changes on v1.20

- `PodPreset`: removed
- `ServiceSpec`: `ipFamily` replaced by `ipFamilies` and `ipFamilyPolicy`
- `ServiceSpec`: `clusterIPs` added
- `ServiceSpec`: `allocateLoadBalancerNodePorts` added
- `HorizontalPodAutoscalerSpec`: `metrics.containerResource` added
- `HorizontalPodAutoscalerStatus`: `currentMetrics.containerResource` added
- `EndpointSlice`: `endpoints.conditions.serving`, `endpoints.conditions.terminating` added
- `FlowSchema` from `v1alpha1` to `v1beta1`
- `PriorityLevelConfiguration` from `v1alpha1` to `v1beta1`
- `node.k8s.io.RuntimeClass` from `v1beta1` to `v1`
