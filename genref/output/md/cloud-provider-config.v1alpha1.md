---
title: Cloud Provider Configuration (v1alpha1)
content_type: tool-reference
package: cloudcontrollermanager.config.k8s.io/v1alpha1
auto_generated: true
---


## Resource Types 


- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)
  
    

## `CloudControllerManagerConfiguration`     {#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration}
    



<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>cloudcontrollermanager.config.k8s.io/v1alpha1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>CloudControllerManagerConfiguration</code></td></tr>
    
  
<tr><td><code>Generic</code> <B>[Required]</B><br/>
<a href="#controllermanager-config-k8s-io-v1alpha1-GenericControllerManagerConfiguration"><code>GenericControllerManagerConfiguration</code></a>
</td>
<td>
   <p>Generic holds configuration for a generic controller-manager</p>
</td>
</tr>
<tr><td><code>KubeCloudShared</code> <B>[Required]</B><br/>
<a href="#cloudcontrollermanager-config-k8s-io-v1alpha1-KubeCloudSharedConfiguration"><code>KubeCloudSharedConfiguration</code></a>
</td>
<td>
   <p>KubeCloudSharedConfiguration holds configuration for shared related features
both in cloud controller manager and kube-controller manager.</p>
</td>
</tr>
<tr><td><code>ServiceController</code> <B>[Required]</B><br/>
<a href="#ServiceControllerConfiguration"><code>ServiceControllerConfiguration</code></a>
</td>
<td>
   <p>ServiceControllerConfiguration holds configuration for ServiceController
related features.</p>
</td>
</tr>
<tr><td><code>NodeStatusUpdateFrequency</code> <B>[Required]</B><br/>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   <p>NodeStatusUpdateFrequency is the frequency at which the controller updates nodes' status</p>
</td>
</tr>
</tbody>
</table>

## `CloudProviderConfiguration`     {#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudProviderConfiguration}
    

**Appears in:**

- [KubeCloudSharedConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-KubeCloudSharedConfiguration)


<p>CloudProviderConfiguration contains basically elements about cloud provider.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>Name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>Name is the provider for cloud services.</p>
</td>
</tr>
<tr><td><code>CloudConfigFile</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>cloudConfigFile is the path to the cloud provider configuration file.</p>
</td>
</tr>
</tbody>
</table>

## `KubeCloudSharedConfiguration`     {#cloudcontrollermanager-config-k8s-io-v1alpha1-KubeCloudSharedConfiguration}
    

**Appears in:**

- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


<p>KubeCloudSharedConfiguration contains elements shared by both kube-controller manager
and cloud-controller manager, but not genericconfig.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>CloudProvider</code> <B>[Required]</B><br/>
<a href="#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudProviderConfiguration"><code>CloudProviderConfiguration</code></a>
</td>
<td>
   <p>CloudProviderConfiguration holds configuration for CloudProvider related features.</p>
</td>
</tr>
<tr><td><code>ExternalCloudVolumePlugin</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>externalCloudVolumePlugin specifies the plugin to use when cloudProvider is &quot;external&quot;.
It is currently used by the in repo cloud providers to handle node and volume control in the KCM.</p>
</td>
</tr>
<tr><td><code>UseServiceAccountCredentials</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   <p>useServiceAccountCredentials indicates whether controllers should be run with
individual service account credentials.</p>
</td>
</tr>
<tr><td><code>AllowUntaggedCloud</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   <p>run with untagged cloud instances</p>
</td>
</tr>
<tr><td><code>RouteReconciliationPeriod</code> <B>[Required]</B><br/>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   <p>routeReconciliationPeriod is the period for reconciling routes created for Nodes by cloud provider..</p>
</td>
</tr>
<tr><td><code>NodeMonitorPeriod</code> <B>[Required]</B><br/>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   <p>nodeMonitorPeriod is the period for syncing NodeStatus in NodeController.</p>
</td>
</tr>
<tr><td><code>ClusterName</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>clusterName is the instance prefix for the cluster.</p>
</td>
</tr>
<tr><td><code>ClusterCIDR</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>clusterCIDR is CIDR Range for Pods in cluster.</p>
</td>
</tr>
<tr><td><code>AllocateNodeCIDRs</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   <p>AllocateNodeCIDRs enables CIDRs for Pods to be allocated and, if
ConfigureCloudRoutes is true, to be set on the cloud provider.</p>
</td>
</tr>
<tr><td><code>CIDRAllocatorType</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>CIDRAllocatorType determines what kind of pod CIDR allocator will be used.</p>
</td>
</tr>
<tr><td><code>ConfigureCloudRoutes</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   <p>configureCloudRoutes enables CIDRs allocated with allocateNodeCIDRs
to be configured on the cloud provider.</p>
</td>
</tr>
<tr><td><code>NodeSyncPeriod</code> <B>[Required]</B><br/>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   <p>nodeSyncPeriod is the period for syncing nodes from cloudprovider. Longer
periods will result in fewer calls to cloud provider, but may delay addition
of new nodes to cluster.</p>
</td>
</tr>
</tbody>
</table>
  
  
    

## `ControllerLeaderConfiguration`     {#controllermanager-config-k8s-io-v1alpha1-ControllerLeaderConfiguration}
    

**Appears in:**

- [LeaderMigrationConfiguration](#controllermanager-config-k8s-io-v1alpha1-LeaderMigrationConfiguration)


<p>ControllerLeaderConfiguration provides the configuration for a migrating leader lock.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>Name is the name of the controller being migrated
E.g. service-controller, route-controller, cloud-node-controller, etc</p>
</td>
</tr>
<tr><td><code>component</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>Component is the name of the component in which the controller should be running.
E.g. kube-controller-manager, cloud-controller-manager, etc
Or '&lowast;' meaning the controller can be run under any component that participates in the migration</p>
</td>
</tr>
</tbody>
</table>

## `GenericControllerManagerConfiguration`     {#controllermanager-config-k8s-io-v1alpha1-GenericControllerManagerConfiguration}
    

**Appears in:**

- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


<p>GenericControllerManagerConfiguration holds configuration for a generic controller-manager.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>Port</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   <p>port is the port that the controller-manager's http service runs on.</p>
</td>
</tr>
<tr><td><code>Address</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>address is the IP address to serve on (set to 0.0.0.0 for all interfaces).</p>
</td>
</tr>
<tr><td><code>MinResyncPeriod</code> <B>[Required]</B><br/>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   <p>minResyncPeriod is the resync period in reflectors; will be random between
minResyncPeriod and 2&lowast;minResyncPeriod.</p>
</td>
</tr>
<tr><td><code>ClientConnection</code> <B>[Required]</B><br/>
<a href="#ClientConnectionConfiguration"><code>ClientConnectionConfiguration</code></a>
</td>
<td>
   <p>ClientConnection specifies the kubeconfig file and client connection
settings for the proxy server to use when communicating with the apiserver.</p>
</td>
</tr>
<tr><td><code>ControllerStartInterval</code> <B>[Required]</B><br/>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   <p>How long to wait between starting controller managers</p>
</td>
</tr>
<tr><td><code>LeaderElection</code> <B>[Required]</B><br/>
<a href="#LeaderElectionConfiguration"><code>LeaderElectionConfiguration</code></a>
</td>
<td>
   <p>leaderElection defines the configuration of leader election client.</p>
</td>
</tr>
<tr><td><code>Controllers</code> <B>[Required]</B><br/>
<code>[]string</code>
</td>
<td>
   <p>Controllers is the list of controllers to enable or disable
'&lowast;' means &quot;all enabled by default controllers&quot;
'foo' means &quot;enable 'foo'&quot;
'-foo' means &quot;disable 'foo'&quot;
first item for a particular name wins</p>
</td>
</tr>
<tr><td><code>Debugging</code> <B>[Required]</B><br/>
<a href="#DebuggingConfiguration"><code>DebuggingConfiguration</code></a>
</td>
<td>
   <p>DebuggingConfiguration holds configuration for Debugging related features.</p>
</td>
</tr>
<tr><td><code>LeaderMigrationEnabled</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   <p>LeaderMigrationEnabled indicates whether Leader Migration should be enabled for the controller manager.</p>
</td>
</tr>
<tr><td><code>LeaderMigration</code> <B>[Required]</B><br/>
<a href="#controllermanager-config-k8s-io-v1alpha1-LeaderMigrationConfiguration"><code>LeaderMigrationConfiguration</code></a>
</td>
<td>
   <p>LeaderMigration holds the configuration for Leader Migration.</p>
</td>
</tr>
</tbody>
</table>

## `LeaderMigrationConfiguration`     {#controllermanager-config-k8s-io-v1alpha1-LeaderMigrationConfiguration}
    

**Appears in:**

- [GenericControllerManagerConfiguration](#controllermanager-config-k8s-io-v1alpha1-GenericControllerManagerConfiguration)


<p>LeaderMigrationConfiguration provides versioned configuration for all migrating leader locks.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>leaderName</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>LeaderName is the name of the leader election resource that protects the migration
E.g. 1-20-KCM-to-1-21-CCM</p>
</td>
</tr>
<tr><td><code>resourceLock</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>ResourceLock indicates the resource object type that will be used to lock
Should be &quot;leases&quot; or &quot;endpoints&quot;</p>
</td>
</tr>
<tr><td><code>controllerLeaders</code> <B>[Required]</B><br/>
<a href="#controllermanager-config-k8s-io-v1alpha1-ControllerLeaderConfiguration"><code>[]ControllerLeaderConfiguration</code></a>
</td>
<td>
   <p>ControllerLeaders contains a list of migrating leader lock configurations</p>
</td>
</tr>
</tbody>
</table>
  
  
    

## `ServiceControllerConfiguration`     {#ServiceControllerConfiguration}
    

**Appears in:**

- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


<p>ServiceControllerConfiguration contains elements describing ServiceController.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>ConcurrentServiceSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   <p>concurrentServiceSyncs is the number of services that are
allowed to sync concurrently. Larger number = more responsive service
management, but more CPU (and network) load.</p>
</td>
</tr>
</tbody>
</table>
