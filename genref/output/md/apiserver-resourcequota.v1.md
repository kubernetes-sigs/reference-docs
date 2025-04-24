---
title: kube-apiserver ResourceQuota Configuration (v1)
content_type: tool-reference
package: apiserver.config.k8s.io/v1
auto_generated: true
---
<p>Package v1 is the v1 version of the API.</p>


## Resource Types 


- [Configuration](#apiserver-config-k8s-io-v1-Configuration)
  

## `Configuration`     {#apiserver-config-k8s-io-v1-Configuration}
    


<p>Configuration provides configuration for the ResourceQuota admission controller.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>apiserver.config.k8s.io/v1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>Configuration</code></td></tr>
    
  
<tr><td><code>limitedResources</code><br/>
<a href="#apiserver-config-k8s-io-v1-LimitedResource"><code>[]LimitedResource</code></a>
</td>
<td>
   <p>LimitedResources whose consumption is limited by default.</p>
</td>
</tr>
</tbody>
</table>

## `LimitedResource`     {#apiserver-config-k8s-io-v1-LimitedResource}
    

**Appears in:**

- [Configuration](#apiserver-config-k8s-io-v1-Configuration)


<p>LimitedResource matches a resource whose consumption is limited by default.
To consume the resource, there must exist an associated quota that limits
its consumption.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>apiGroup</code><br/>
<code>string</code>
</td>
<td>
   <p>APIGroup is the name of the APIGroup that contains the limited resource.</p>
</td>
</tr>
<tr><td><code>resource</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>Resource is the name of the resource this rule applies to.
For example, if the administrator wants to limit consumption
of a storage resource associated with persistent volume claims,
the value would be &quot;persistentvolumeclaims&quot;.</p>
</td>
</tr>
<tr><td><code>matchContains</code> <B>[Required]</B><br/>
<code>[]string</code>
</td>
<td>
   <p>For each intercepted request, the quota system will evaluate
its resource usage.  It will iterate through each resource consumed
and if the resource contains any substring in this listing, the
quota system will ensure that there is a covering quota.  In the
absence of a covering quota, the quota system will deny the request.
For example, if an administrator wants to globally enforce that
that a quota must exist to consume persistent volume claims associated
with any storage class, the list would include
&quot;.storageclass.storage.k8s.io/requests.storage&quot;</p>
</td>
</tr>
<tr><td><code>matchScopes</code><br/>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.33/#scopedresourceselectorrequirement-v1-core"><code>[]core/v1.ScopedResourceSelectorRequirement</code></a>
</td>
<td>
   <p>For each intercepted request, the quota system will figure out if the input object
satisfies a scope which is present in this listing, then
quota system will ensure that there is a covering quota.  In the
absence of a covering quota, the quota system will deny the request.
For example, if an administrator wants to globally enforce that
a quota must exist to create a pod with &quot;cluster-services&quot; priorityclass
the list would include &quot;scopeName=PriorityClass, Operator=In, Value=cluster-services&quot;</p>
</td>
</tr>
</tbody>
</table>
  