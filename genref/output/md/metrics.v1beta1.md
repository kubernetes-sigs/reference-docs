---
title: Kubernetes Metrics (v1alpha1)
content_type: tool-reference
package: metrics.k8s.io/v1beta1
auto_generated: true
---
Package v1beta1 is the v1beta1 version of the metrics API.

## Resource Types 


- [NodeMetrics](#metrics-k8s-io-v1beta1-NodeMetrics)
- [NodeMetricsList](#metrics-k8s-io-v1beta1-NodeMetricsList)
- [PodMetrics](#metrics-k8s-io-v1beta1-PodMetrics)
- [PodMetricsList](#metrics-k8s-io-v1beta1-PodMetricsList)
  
    


## `NodeMetrics`     {#metrics-k8s-io-v1beta1-NodeMetrics}
    



**Appears in:**

- [NodeMetricsList](#metrics-k8s-io-v1beta1-NodeMetricsList)


NodeMetrics sets resource usage metrics of a node.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>metrics.k8s.io/v1beta1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>NodeMetrics</code></td></tr>
    

  
  
<tr><td><code>metadata</code> <B>[Required]</B><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta"><code>meta/v1.ObjectMeta</code></a>
</td>
<td>
   <span class="text-muted">No description provided.</span>
   Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
    
  
<tr><td><code>timestamp</code> <B>[Required]</B><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#time-v1-meta"><code>meta/v1.Time</code></a>
</td>
<td>
   The following fields define time interval from which metrics were
collected from the interval [Timestamp-Window, Timestamp].</td>
</tr>
    
  
<tr><td><code>window</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   <span class="text-muted">No description provided.</span>
   </td>
</tr>
    
  
<tr><td><code>usage</code> <B>[Required]</B><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#resourcelist-v1-core"><code>core/v1.ResourceList</code></a>
</td>
<td>
   The memory usage is the memory working set.</td>
</tr>
    
  
</tbody>
</table>
    


## `NodeMetricsList`     {#metrics-k8s-io-v1beta1-NodeMetricsList}
    




NodeMetricsList is a list of NodeMetrics.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>metrics.k8s.io/v1beta1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>NodeMetricsList</code></td></tr>
    

  
  
<tr><td><code>metadata</code> <B>[Required]</B><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#listmeta-v1-meta"><code>meta/v1.ListMeta</code></a>
</td>
<td>
   Standard list metadata.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds</td>
</tr>
    
  
<tr><td><code>items</code> <B>[Required]</B><br/>
<a href="#metrics-k8s-io-v1beta1-NodeMetrics"><code>[]NodeMetrics</code></a>
</td>
<td>
   List of node metrics.</td>
</tr>
    
  
</tbody>
</table>
    


## `PodMetrics`     {#metrics-k8s-io-v1beta1-PodMetrics}
    



**Appears in:**

- [PodMetricsList](#metrics-k8s-io-v1beta1-PodMetricsList)


PodMetrics sets resource usage metrics of a pod.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>metrics.k8s.io/v1beta1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>PodMetrics</code></td></tr>
    

  
  
<tr><td><code>metadata</code> <B>[Required]</B><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta"><code>meta/v1.ObjectMeta</code></a>
</td>
<td>
   <span class="text-muted">No description provided.</span>
   Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
    
  
<tr><td><code>timestamp</code> <B>[Required]</B><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#time-v1-meta"><code>meta/v1.Time</code></a>
</td>
<td>
   The following fields define time interval from which metrics were
collected from the interval [Timestamp-Window, Timestamp].</td>
</tr>
    
  
<tr><td><code>window</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   <span class="text-muted">No description provided.</span>
   </td>
</tr>
    
  
<tr><td><code>containers</code> <B>[Required]</B><br/>
<a href="#metrics-k8s-io-v1beta1-ContainerMetrics"><code>[]ContainerMetrics</code></a>
</td>
<td>
   Metrics for all containers are collected within the same time window.</td>
</tr>
    
  
</tbody>
</table>
    


## `PodMetricsList`     {#metrics-k8s-io-v1beta1-PodMetricsList}
    




PodMetricsList is a list of PodMetrics.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>metrics.k8s.io/v1beta1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>PodMetricsList</code></td></tr>
    

  
  
<tr><td><code>metadata</code> <B>[Required]</B><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#listmeta-v1-meta"><code>meta/v1.ListMeta</code></a>
</td>
<td>
   Standard list metadata.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds</td>
</tr>
    
  
<tr><td><code>items</code> <B>[Required]</B><br/>
<a href="#metrics-k8s-io-v1beta1-PodMetrics"><code>[]PodMetrics</code></a>
</td>
<td>
   List of pod metrics.</td>
</tr>
    
  
</tbody>
</table>
    


## `ContainerMetrics`     {#metrics-k8s-io-v1beta1-ContainerMetrics}
    



**Appears in:**

- [PodMetrics](#metrics-k8s-io-v1beta1-PodMetrics)


ContainerMetrics sets resource usage metrics of a container.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   Container name corresponding to the one from pod.spec.containers.</td>
</tr>
    
  
<tr><td><code>usage</code> <B>[Required]</B><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#resourcelist-v1-core"><code>core/v1.ResourceList</code></a>
</td>
<td>
   The memory usage is the memory working set.</td>
</tr>
    
  
</tbody>
</table>
    
  
