---
title: kube-apiserver ResourceQuota Configuration (v1)
content_type: tool-reference
package: apiserver.config.k8s.io/v1
auto_generated: true
---
Package v1 is the v1 version of the API.

## Resource Types 


- [Configuration](#apiserver-config-k8s-io-v1-Configuration)
  
    


## `Configuration`     {#apiserver-config-k8s-io-v1-Configuration}
    




Configuration provides configuration for the ResourceQuota admission controller.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>apiserver.config.k8s.io/v1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>Configuration</code></td></tr>
    

  
  
<tr><td><code>limitedResources</code><br/>
<a href="#apiserver-config-k8s-io-v1-LimitedResource"><code>[]LimitedResource</code></a>
</td>
<td>
   LimitedResources whose consumption is limited by default.</td>
</tr>
    
  
</tbody>
</table>
    


## `LimitedResource`     {#apiserver-config-k8s-io-v1-LimitedResource}
    



**Appears in:**

- [Configuration](#apiserver-config-k8s-io-v1-Configuration)


LimitedResource matches a resource whose consumption is limited by default.
To consume the resource, there must exist an associated quota that limits
its consumption.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>apiGroup</code><br/>
<code>string</code>
</td>
<td>
   APIGroup is the name of the APIGroup that contains the limited resource.</td>
</tr>
    
  
<tr><td><code>resource</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   Resource is the name of the resource this rule applies to.
For example, if the administrator wants to limit consumption
of a storage resource associated with persistent volume claims,
the value would be "persistentvolumeclaims".</td>
</tr>
    
  
<tr><td><code>matchContains</code> <B>[Required]</B><br/>
<code>[]string</code>
</td>
<td>
   For each intercepted request, the quota system will evaluate
its resource usage.  It will iterate through each resource consumed
and if the resource contains any substring in this listing, the
quota system will ensure that there is a covering quota.  In the
absence of a covering quota, the quota system will deny the request.
For example, if an administrator wants to globally enforce that
that a quota must exist to consume persistent volume claims associated
with any storage class, the list would include
".storageclass.storage.k8s.io/requests.storage"</td>
</tr>
    
  
<tr><td><code>matchScopes</code><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#scopedresourceselectorrequirement-v1-core"><code>[]core/v1.ScopedResourceSelectorRequirement</code></a>
</td>
<td>
   For each intercepted request, the quota system will figure out if the input object
satisfies a scope which is present in this listing, then
quota system will ensure that there is a covering quota.  In the
absence of a covering quota, the quota system will deny the request.
For example, if an administrator wants to globally enforce that
a quota must exist to create a pod with "cluster-services" priorityclass
the list would include "scopeName=PriorityClass, Operator=In, Value=cluster-services"</td>
</tr>
    
  
</tbody>
</table>
    
  
