

## Resource Types 


  
    
    

## `NodeControllerConfiguration`     {#NodeControllerConfiguration}
    

**Appears in:**

- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)


<p>NodeControllerConfiguration contains elements describing NodeController.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>ConcurrentNodeSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   <p>ConcurrentNodeSyncs is the number of workers
concurrently synchronizing nodes</p>
</td>
</tr>
<tr><td><code>ConcurrentNodeStatusUpdates</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   <p>ConcurrentNodeStatusUpdates is the number of workers
concurrently updating node statuses.
If unspecified or 0, ConcurrentNodeSyncs is used instead</p>
</td>
</tr>
</tbody>
</table>