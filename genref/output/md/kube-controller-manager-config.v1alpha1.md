---
title: kube-controller-manager Configuration (v1alpha1)
content_type: tool-reference
package: cloudcontrollermanager.config.k8s.io/v1alpha1
auto_generated: true
---


## Resource Types 


- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)
- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)
  
    


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
   Generic holds configuration for a generic controller-manager</td>
</tr>
    
  
<tr><td><code>KubeCloudShared</code> <B>[Required]</B><br/>
<a href="#cloudcontrollermanager-config-k8s-io-v1alpha1-KubeCloudSharedConfiguration"><code>KubeCloudSharedConfiguration</code></a>
</td>
<td>
   KubeCloudSharedConfiguration holds configuration for shared related features
both in cloud controller manager and kube-controller manager.</td>
</tr>
    
  
<tr><td><code>ServiceController</code> <B>[Required]</B><br/>
<a href="#ServiceControllerConfiguration"><code>ServiceControllerConfiguration</code></a>
</td>
<td>
   ServiceControllerConfiguration holds configuration for ServiceController
related features.</td>
</tr>
    
  
<tr><td><code>NodeStatusUpdateFrequency</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   NodeStatusUpdateFrequency is the frequency at which the controller updates nodes' status</td>
</tr>
    
  
</tbody>
</table>
    


## `CloudProviderConfiguration`     {#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudProviderConfiguration}
    



**Appears in:**

- [KubeCloudSharedConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-KubeCloudSharedConfiguration)


CloudProviderConfiguration contains basically elements about cloud provider.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>Name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   Name is the provider for cloud services.</td>
</tr>
    
  
<tr><td><code>CloudConfigFile</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   cloudConfigFile is the path to the cloud provider configuration file.</td>
</tr>
    
  
</tbody>
</table>
    


## `KubeCloudSharedConfiguration`     {#cloudcontrollermanager-config-k8s-io-v1alpha1-KubeCloudSharedConfiguration}
    



**Appears in:**

- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


KubeCloudSharedConfiguration contains elements shared by both kube-controller manager
and cloud-controller manager, but not genericconfig.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>CloudProvider</code> <B>[Required]</B><br/>
<a href="#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudProviderConfiguration"><code>CloudProviderConfiguration</code></a>
</td>
<td>
   CloudProviderConfiguration holds configuration for CloudProvider related features.</td>
</tr>
    
  
<tr><td><code>ExternalCloudVolumePlugin</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   externalCloudVolumePlugin specifies the plugin to use when cloudProvider is "external".
It is currently used by the in repo cloud providers to handle node and volume control in the KCM.</td>
</tr>
    
  
<tr><td><code>UseServiceAccountCredentials</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   useServiceAccountCredentials indicates whether controllers should be run with
individual service account credentials.</td>
</tr>
    
  
<tr><td><code>AllowUntaggedCloud</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   run with untagged cloud instances</td>
</tr>
    
  
<tr><td><code>RouteReconciliationPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   routeReconciliationPeriod is the period for reconciling routes created for Nodes by cloud provider..</td>
</tr>
    
  
<tr><td><code>NodeMonitorPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   nodeMonitorPeriod is the period for syncing NodeStatus in NodeController.</td>
</tr>
    
  
<tr><td><code>ClusterName</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   clusterName is the instance prefix for the cluster.</td>
</tr>
    
  
<tr><td><code>ClusterCIDR</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   clusterCIDR is CIDR Range for Pods in cluster.</td>
</tr>
    
  
<tr><td><code>AllocateNodeCIDRs</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   AllocateNodeCIDRs enables CIDRs for Pods to be allocated and, if
ConfigureCloudRoutes is true, to be set on the cloud provider.</td>
</tr>
    
  
<tr><td><code>CIDRAllocatorType</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   CIDRAllocatorType determines what kind of pod CIDR allocator will be used.</td>
</tr>
    
  
<tr><td><code>ConfigureCloudRoutes</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   configureCloudRoutes enables CIDRs allocated with allocateNodeCIDRs
to be configured on the cloud provider.</td>
</tr>
    
  
<tr><td><code>NodeSyncPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   nodeSyncPeriod is the period for syncing nodes from cloudprovider. Longer
periods will result in fewer calls to cloud provider, but may delay addition
of new nodes to cluster.</td>
</tr>
    
  
</tbody>
</table>
    
  
  
    


## `ControllerLeaderConfiguration`     {#controllermanager-config-k8s-io-v1alpha1-ControllerLeaderConfiguration}
    



**Appears in:**

- [LeaderMigrationConfiguration](#controllermanager-config-k8s-io-v1alpha1-LeaderMigrationConfiguration)


ControllerLeaderConfiguration provides the configuration for a migrating leader lock.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   Name is the name of the controller being migrated
E.g. service-controller, route-controller, cloud-node-controller, etc</td>
</tr>
    
  
<tr><td><code>component</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   Component is the name of the component in which the controller should be running.
E.g. kube-controller-manager, cloud-controller-manager, etc
Or '&lowast;' meaning the controller can be run under any component that participates in the migration</td>
</tr>
    
  
</tbody>
</table>
    


## `GenericControllerManagerConfiguration`     {#controllermanager-config-k8s-io-v1alpha1-GenericControllerManagerConfiguration}
    



**Appears in:**

- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


GenericControllerManagerConfiguration holds configuration for a generic controller-manager.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>Port</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   port is the port that the controller-manager's http service runs on.</td>
</tr>
    
  
<tr><td><code>Address</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   address is the IP address to serve on (set to 0.0.0.0 for all interfaces).</td>
</tr>
    
  
<tr><td><code>MinResyncPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   minResyncPeriod is the resync period in reflectors; will be random between
minResyncPeriod and 2&lowast;minResyncPeriod.</td>
</tr>
    
  
<tr><td><code>ClientConnection</code> <B>[Required]</B><br/>
<a href="#ClientConnectionConfiguration"><code>ClientConnectionConfiguration</code></a>
</td>
<td>
   ClientConnection specifies the kubeconfig file and client connection
settings for the proxy server to use when communicating with the apiserver.</td>
</tr>
    
  
<tr><td><code>ControllerStartInterval</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   How long to wait between starting controller managers</td>
</tr>
    
  
<tr><td><code>LeaderElection</code> <B>[Required]</B><br/>
<a href="#LeaderElectionConfiguration"><code>LeaderElectionConfiguration</code></a>
</td>
<td>
   leaderElection defines the configuration of leader election client.</td>
</tr>
    
  
<tr><td><code>Controllers</code> <B>[Required]</B><br/>
<code>[]string</code>
</td>
<td>
   Controllers is the list of controllers to enable or disable
'&lowast;' means "all enabled by default controllers"
'foo' means "enable 'foo'"
'-foo' means "disable 'foo'"
first item for a particular name wins</td>
</tr>
    
  
<tr><td><code>Debugging</code> <B>[Required]</B><br/>
<a href="#DebuggingConfiguration"><code>DebuggingConfiguration</code></a>
</td>
<td>
   DebuggingConfiguration holds configuration for Debugging related features.</td>
</tr>
    
  
<tr><td><code>LeaderMigrationEnabled</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   LeaderMigrationEnabled indicates whether Leader Migration should be enabled for the controller manager.</td>
</tr>
    
  
<tr><td><code>LeaderMigration</code> <B>[Required]</B><br/>
<a href="#controllermanager-config-k8s-io-v1alpha1-LeaderMigrationConfiguration"><code>LeaderMigrationConfiguration</code></a>
</td>
<td>
   LeaderMigration holds the configuration for Leader Migration.</td>
</tr>
    
  
</tbody>
</table>
    


## `LeaderMigrationConfiguration`     {#controllermanager-config-k8s-io-v1alpha1-LeaderMigrationConfiguration}
    



**Appears in:**

- [GenericControllerManagerConfiguration](#controllermanager-config-k8s-io-v1alpha1-GenericControllerManagerConfiguration)


LeaderMigrationConfiguration provides versioned configuration for all migrating leader locks.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
  
<tr><td><code>leaderName</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   LeaderName is the name of the leader election resource that protects the migration
E.g. 1-20-KCM-to-1-21-CCM</td>
</tr>
    
  
<tr><td><code>resourceLock</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   ResourceLock indicates the resource object type that will be used to lock
Should be "leases" or "endpoints"</td>
</tr>
    
  
<tr><td><code>controllerLeaders</code> <B>[Required]</B><br/>
<a href="#controllermanager-config-k8s-io-v1alpha1-ControllerLeaderConfiguration"><code>[]ControllerLeaderConfiguration</code></a>
</td>
<td>
   ControllerLeaders contains a list of migrating leader lock configurations</td>
</tr>
    
  
</tbody>
</table>
    
  
  
    


## `KubeControllerManagerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration}
    




KubeControllerManagerConfiguration contains elements describing kube-controller manager.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>kubecontrollermanager.config.k8s.io/v1alpha1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>KubeControllerManagerConfiguration</code></td></tr>
    

  
  
<tr><td><code>Generic</code> <B>[Required]</B><br/>
<a href="#controllermanager-config-k8s-io-v1alpha1-GenericControllerManagerConfiguration"><code>GenericControllerManagerConfiguration</code></a>
</td>
<td>
   Generic holds configuration for a generic controller-manager</td>
</tr>
    
  
<tr><td><code>KubeCloudShared</code> <B>[Required]</B><br/>
<a href="#cloudcontrollermanager-config-k8s-io-v1alpha1-KubeCloudSharedConfiguration"><code>KubeCloudSharedConfiguration</code></a>
</td>
<td>
   KubeCloudSharedConfiguration holds configuration for shared related features
both in cloud controller manager and kube-controller manager.</td>
</tr>
    
  
<tr><td><code>AttachDetachController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-AttachDetachControllerConfiguration"><code>AttachDetachControllerConfiguration</code></a>
</td>
<td>
   AttachDetachControllerConfiguration holds configuration for
AttachDetachController related features.</td>
</tr>
    
  
<tr><td><code>CSRSigningController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-CSRSigningControllerConfiguration"><code>CSRSigningControllerConfiguration</code></a>
</td>
<td>
   CSRSigningControllerConfiguration holds configuration for
CSRSigningController related features.</td>
</tr>
    
  
<tr><td><code>DaemonSetController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-DaemonSetControllerConfiguration"><code>DaemonSetControllerConfiguration</code></a>
</td>
<td>
   DaemonSetControllerConfiguration holds configuration for DaemonSetController
related features.</td>
</tr>
    
  
<tr><td><code>DeploymentController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-DeploymentControllerConfiguration"><code>DeploymentControllerConfiguration</code></a>
</td>
<td>
   DeploymentControllerConfiguration holds configuration for
DeploymentController related features.</td>
</tr>
    
  
<tr><td><code>StatefulSetController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-StatefulSetControllerConfiguration"><code>StatefulSetControllerConfiguration</code></a>
</td>
<td>
   StatefulSetControllerConfiguration holds configuration for
StatefulSetController related features.</td>
</tr>
    
  
<tr><td><code>DeprecatedController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-DeprecatedControllerConfiguration"><code>DeprecatedControllerConfiguration</code></a>
</td>
<td>
   DeprecatedControllerConfiguration holds configuration for some deprecated
features.</td>
</tr>
    
  
<tr><td><code>EndpointController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-EndpointControllerConfiguration"><code>EndpointControllerConfiguration</code></a>
</td>
<td>
   EndpointControllerConfiguration holds configuration for EndpointController
related features.</td>
</tr>
    
  
<tr><td><code>EndpointSliceController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-EndpointSliceControllerConfiguration"><code>EndpointSliceControllerConfiguration</code></a>
</td>
<td>
   EndpointSliceControllerConfiguration holds configuration for
EndpointSliceController related features.</td>
</tr>
    
  
<tr><td><code>EndpointSliceMirroringController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-EndpointSliceMirroringControllerConfiguration"><code>EndpointSliceMirroringControllerConfiguration</code></a>
</td>
<td>
   EndpointSliceMirroringControllerConfiguration holds configuration for
EndpointSliceMirroringController related features.</td>
</tr>
    
  
<tr><td><code>EphemeralVolumeController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-EphemeralVolumeControllerConfiguration"><code>EphemeralVolumeControllerConfiguration</code></a>
</td>
<td>
   EphemeralVolumeControllerConfiguration holds configuration for EphemeralVolumeController
related features.</td>
</tr>
    
  
<tr><td><code>GarbageCollectorController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-GarbageCollectorControllerConfiguration"><code>GarbageCollectorControllerConfiguration</code></a>
</td>
<td>
   GarbageCollectorControllerConfiguration holds configuration for
GarbageCollectorController related features.</td>
</tr>
    
  
<tr><td><code>HPAController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-HPAControllerConfiguration"><code>HPAControllerConfiguration</code></a>
</td>
<td>
   HPAControllerConfiguration holds configuration for HPAController related features.</td>
</tr>
    
  
<tr><td><code>JobController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-JobControllerConfiguration"><code>JobControllerConfiguration</code></a>
</td>
<td>
   JobControllerConfiguration holds configuration for JobController related features.</td>
</tr>
    
  
<tr><td><code>CronJobController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-CronJobControllerConfiguration"><code>CronJobControllerConfiguration</code></a>
</td>
<td>
   CronJobControllerConfiguration holds configuration for CronJobController related features.</td>
</tr>
    
  
<tr><td><code>NamespaceController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-NamespaceControllerConfiguration"><code>NamespaceControllerConfiguration</code></a>
</td>
<td>
   NamespaceControllerConfiguration holds configuration for NamespaceController
related features.
NamespaceControllerConfiguration holds configuration for NamespaceController
related features.</td>
</tr>
    
  
<tr><td><code>NodeIPAMController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-NodeIPAMControllerConfiguration"><code>NodeIPAMControllerConfiguration</code></a>
</td>
<td>
   NodeIPAMControllerConfiguration holds configuration for NodeIPAMController
related features.</td>
</tr>
    
  
<tr><td><code>NodeLifecycleController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-NodeLifecycleControllerConfiguration"><code>NodeLifecycleControllerConfiguration</code></a>
</td>
<td>
   NodeLifecycleControllerConfiguration holds configuration for
NodeLifecycleController related features.</td>
</tr>
    
  
<tr><td><code>PersistentVolumeBinderController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-PersistentVolumeBinderControllerConfiguration"><code>PersistentVolumeBinderControllerConfiguration</code></a>
</td>
<td>
   PersistentVolumeBinderControllerConfiguration holds configuration for
PersistentVolumeBinderController related features.</td>
</tr>
    
  
<tr><td><code>PodGCController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-PodGCControllerConfiguration"><code>PodGCControllerConfiguration</code></a>
</td>
<td>
   PodGCControllerConfiguration holds configuration for PodGCController
related features.</td>
</tr>
    
  
<tr><td><code>ReplicaSetController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-ReplicaSetControllerConfiguration"><code>ReplicaSetControllerConfiguration</code></a>
</td>
<td>
   ReplicaSetControllerConfiguration holds configuration for ReplicaSet related features.</td>
</tr>
    
  
<tr><td><code>ReplicationController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-ReplicationControllerConfiguration"><code>ReplicationControllerConfiguration</code></a>
</td>
<td>
   ReplicationControllerConfiguration holds configuration for
ReplicationController related features.</td>
</tr>
    
  
<tr><td><code>ResourceQuotaController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-ResourceQuotaControllerConfiguration"><code>ResourceQuotaControllerConfiguration</code></a>
</td>
<td>
   ResourceQuotaControllerConfiguration holds configuration for
ResourceQuotaController related features.</td>
</tr>
    
  
<tr><td><code>SAController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-SAControllerConfiguration"><code>SAControllerConfiguration</code></a>
</td>
<td>
   SAControllerConfiguration holds configuration for ServiceAccountController
related features.</td>
</tr>
    
  
<tr><td><code>ServiceController</code> <B>[Required]</B><br/>
<a href="#ServiceControllerConfiguration"><code>ServiceControllerConfiguration</code></a>
</td>
<td>
   ServiceControllerConfiguration holds configuration for ServiceController
related features.</td>
</tr>
    
  
<tr><td><code>TTLAfterFinishedController</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-TTLAfterFinishedControllerConfiguration"><code>TTLAfterFinishedControllerConfiguration</code></a>
</td>
<td>
   TTLAfterFinishedControllerConfiguration holds configuration for
TTLAfterFinishedController related features.</td>
</tr>
    
  
</tbody>
</table>
    


## `AttachDetachControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-AttachDetachControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


AttachDetachControllerConfiguration contains elements describing AttachDetachController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>DisableAttachDetachReconcilerSync</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   Reconciler runs a periodic loop to reconcile the desired state of the with
the actual state of the world by triggering attach detach operations.
This flag enables or disables reconcile.  Is false by default, and thus enabled.</td>
</tr>
    
  
<tr><td><code>ReconcilerSyncLoopPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   ReconcilerSyncLoopPeriod is the amount of time the reconciler sync states loop
wait between successive executions. Is set to 5 sec by default.</td>
</tr>
    
  
</tbody>
</table>
    


## `CSRSigningConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-CSRSigningConfiguration}
    



**Appears in:**

- [CSRSigningControllerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-CSRSigningControllerConfiguration)


CSRSigningConfiguration holds information about a particular CSR signer

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>CertFile</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   certFile is the filename containing a PEM-encoded
X509 CA certificate used to issue certificates</td>
</tr>
    
  
<tr><td><code>KeyFile</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   keyFile is the filename containing a PEM-encoded
RSA or ECDSA private key used to issue certificates</td>
</tr>
    
  
</tbody>
</table>
    


## `CSRSigningControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-CSRSigningControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


CSRSigningControllerConfiguration contains elements describing CSRSigningController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ClusterSigningCertFile</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   clusterSigningCertFile is the filename containing a PEM-encoded
X509 CA certificate used to issue cluster-scoped certificates</td>
</tr>
    
  
<tr><td><code>ClusterSigningKeyFile</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   clusterSigningCertFile is the filename containing a PEM-encoded
RSA or ECDSA private key used to issue cluster-scoped certificates</td>
</tr>
    
  
<tr><td><code>KubeletServingSignerConfiguration</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-CSRSigningConfiguration"><code>CSRSigningConfiguration</code></a>
</td>
<td>
   kubeletServingSignerConfiguration holds the certificate and key used to issue certificates for the kubernetes.io/kubelet-serving signer</td>
</tr>
    
  
<tr><td><code>KubeletClientSignerConfiguration</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-CSRSigningConfiguration"><code>CSRSigningConfiguration</code></a>
</td>
<td>
   kubeletClientSignerConfiguration holds the certificate and key used to issue certificates for the kubernetes.io/kube-apiserver-client-kubelet</td>
</tr>
    
  
<tr><td><code>KubeAPIServerClientSignerConfiguration</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-CSRSigningConfiguration"><code>CSRSigningConfiguration</code></a>
</td>
<td>
   kubeAPIServerClientSignerConfiguration holds the certificate and key used to issue certificates for the kubernetes.io/kube-apiserver-client</td>
</tr>
    
  
<tr><td><code>LegacyUnknownSignerConfiguration</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-CSRSigningConfiguration"><code>CSRSigningConfiguration</code></a>
</td>
<td>
   legacyUnknownSignerConfiguration holds the certificate and key used to issue certificates for the kubernetes.io/legacy-unknown</td>
</tr>
    
  
<tr><td><code>ClusterSigningDuration</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   clusterSigningDuration is the max length of duration signed certificates will be given.
Individual CSRs may request shorter certs by setting spec.expirationSeconds.</td>
</tr>
    
  
</tbody>
</table>
    


## `CronJobControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-CronJobControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


CronJobControllerConfiguration contains elements describing CrongJob2Controller.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentCronJobSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentCronJobSyncs is the number of job objects that are
allowed to sync concurrently. Larger number = more responsive jobs,
but more CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
    


## `DaemonSetControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-DaemonSetControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


DaemonSetControllerConfiguration contains elements describing DaemonSetController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentDaemonSetSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentDaemonSetSyncs is the number of daemonset objects that are
allowed to sync concurrently. Larger number = more responsive daemonset,
but more CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
    


## `DeploymentControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-DeploymentControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


DeploymentControllerConfiguration contains elements describing DeploymentController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentDeploymentSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentDeploymentSyncs is the number of deployment objects that are
allowed to sync concurrently. Larger number = more responsive deployments,
but more CPU (and network) load.</td>
</tr>
    
  
<tr><td><code>DeploymentControllerSyncPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   deploymentControllerSyncPeriod is the period for syncing the deployments.</td>
</tr>
    
  
</tbody>
</table>
    


## `DeprecatedControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-DeprecatedControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


DeprecatedControllerConfiguration contains elements be deprecated.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>DeletingPodsQPS</code> <B>[Required]</B><br/>
<code>float32</code>
</td>
<td>
   DEPRECATED: deletingPodsQps is the number of nodes per second on which pods are deleted in
case of node failure.</td>
</tr>
    
  
<tr><td><code>DeletingPodsBurst</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   DEPRECATED: deletingPodsBurst is the number of nodes on which pods are bursty deleted in
case of node failure. For more details look into RateLimiter.</td>
</tr>
    
  
<tr><td><code>RegisterRetryCount</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   registerRetryCount is the number of retries for initial node registration.
Retry interval equals node-sync-period.</td>
</tr>
    
  
</tbody>
</table>
    


## `EndpointControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-EndpointControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


EndpointControllerConfiguration contains elements describing EndpointController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentEndpointSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentEndpointSyncs is the number of endpoint syncing operations
that will be done concurrently. Larger number = faster endpoint updating,
but more CPU (and network) load.</td>
</tr>
    
  
<tr><td><code>EndpointUpdatesBatchPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   EndpointUpdatesBatchPeriod describes the length of endpoint updates batching period.
Processing of pod changes will be delayed by this duration to join them with potential
upcoming updates and reduce the overall number of endpoints updates.</td>
</tr>
    
  
</tbody>
</table>
    


## `EndpointSliceControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-EndpointSliceControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


EndpointSliceControllerConfiguration contains elements describing
EndpointSliceController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentServiceEndpointSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentServiceEndpointSyncs is the number of service endpoint syncing
operations that will be done concurrently. Larger number = faster
endpoint slice updating, but more CPU (and network) load.</td>
</tr>
    
  
<tr><td><code>MaxEndpointsPerSlice</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   maxEndpointsPerSlice is the maximum number of endpoints that will be
added to an EndpointSlice. More endpoints per slice will result in fewer
and larger endpoint slices, but larger resources.</td>
</tr>
    
  
<tr><td><code>EndpointUpdatesBatchPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   EndpointUpdatesBatchPeriod describes the length of endpoint updates batching period.
Processing of pod changes will be delayed by this duration to join them with potential
upcoming updates and reduce the overall number of endpoints updates.</td>
</tr>
    
  
</tbody>
</table>
    


## `EndpointSliceMirroringControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-EndpointSliceMirroringControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


EndpointSliceMirroringControllerConfiguration contains elements describing
EndpointSliceMirroringController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>MirroringConcurrentServiceEndpointSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   mirroringConcurrentServiceEndpointSyncs is the number of service endpoint
syncing operations that will be done concurrently. Larger number = faster
endpoint slice updating, but more CPU (and network) load.</td>
</tr>
    
  
<tr><td><code>MirroringMaxEndpointsPerSubset</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   mirroringMaxEndpointsPerSubset is the maximum number of endpoints that
will be mirrored to an EndpointSlice for an EndpointSubset.</td>
</tr>
    
  
<tr><td><code>MirroringEndpointUpdatesBatchPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   mirroringEndpointUpdatesBatchPeriod can be used to batch EndpointSlice
updates. All updates triggered by EndpointSlice changes will be delayed
by up to 'mirroringEndpointUpdatesBatchPeriod'. If other addresses in the
same Endpoints resource change in that period, they will be batched to a
single EndpointSlice update. Default 0 value means that each Endpoints
update triggers an EndpointSlice update.</td>
</tr>
    
  
</tbody>
</table>
    


## `EphemeralVolumeControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-EphemeralVolumeControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


EphemeralVolumeControllerConfiguration contains elements describing EphemeralVolumeController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentEphemeralVolumeSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   ConcurrentEphemeralVolumeSyncseSyncs is the number of ephemeral volume syncing operations
that will be done concurrently. Larger number = faster ephemeral volume updating,
but more CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
    


## `GarbageCollectorControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-GarbageCollectorControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


GarbageCollectorControllerConfiguration contains elements describing GarbageCollectorController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>EnableGarbageCollector</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   enables the generic garbage collector. MUST be synced with the
corresponding flag of the kube-apiserver. WARNING: the generic garbage
collector is an alpha feature.</td>
</tr>
    
  
<tr><td><code>ConcurrentGCSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentGCSyncs is the number of garbage collector workers that are
allowed to sync concurrently.</td>
</tr>
    
  
<tr><td><code>GCIgnoredResources</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-GroupResource"><code>[]GroupResource</code></a>
</td>
<td>
   gcIgnoredResources is the list of GroupResources that garbage collection should ignore.</td>
</tr>
    
  
</tbody>
</table>
    


## `GroupResource`     {#kubecontrollermanager-config-k8s-io-v1alpha1-GroupResource}
    



**Appears in:**

- [GarbageCollectorControllerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-GarbageCollectorControllerConfiguration)


GroupResource describes an group resource.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>Group</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   group is the group portion of the GroupResource.</td>
</tr>
    
  
<tr><td><code>Resource</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   resource is the resource portion of the GroupResource.</td>
</tr>
    
  
</tbody>
</table>
    


## `HPAControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-HPAControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


HPAControllerConfiguration contains elements describing HPAController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>HorizontalPodAutoscalerSyncPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   HorizontalPodAutoscalerSyncPeriod is the period for syncing the number of
pods in horizontal pod autoscaler.</td>
</tr>
    
  
<tr><td><code>HorizontalPodAutoscalerUpscaleForbiddenWindow</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   HorizontalPodAutoscalerUpscaleForbiddenWindow is a period after which next upscale allowed.</td>
</tr>
    
  
<tr><td><code>HorizontalPodAutoscalerDownscaleStabilizationWindow</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   HorizontalPodAutoscalerDowncaleStabilizationWindow is a period for which autoscaler will look
backwards and not scale down below any recommendation it made during that period.</td>
</tr>
    
  
<tr><td><code>HorizontalPodAutoscalerDownscaleForbiddenWindow</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   HorizontalPodAutoscalerDownscaleForbiddenWindow is a period after which next downscale allowed.</td>
</tr>
    
  
<tr><td><code>HorizontalPodAutoscalerTolerance</code> <B>[Required]</B><br/>
<code>float64</code>
</td>
<td>
   HorizontalPodAutoscalerTolerance is the tolerance for when
resource usage suggests upscaling/downscaling</td>
</tr>
    
  
<tr><td><code>HorizontalPodAutoscalerCPUInitializationPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   HorizontalPodAutoscalerCPUInitializationPeriod is the period after pod start when CPU samples
might be skipped.</td>
</tr>
    
  
<tr><td><code>HorizontalPodAutoscalerInitialReadinessDelay</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   HorizontalPodAutoscalerInitialReadinessDelay is period after pod start during which readiness
changes are treated as readiness being set for the first time. The only effect of this is that
HPA will disregard CPU samples from unready pods that had last readiness change during that
period.</td>
</tr>
    
  
</tbody>
</table>
    


## `JobControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-JobControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


JobControllerConfiguration contains elements describing JobController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentJobSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentJobSyncs is the number of job objects that are
allowed to sync concurrently. Larger number = more responsive jobs,
but more CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
    


## `NamespaceControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-NamespaceControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


NamespaceControllerConfiguration contains elements describing NamespaceController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>NamespaceSyncPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   namespaceSyncPeriod is the period for syncing namespace life-cycle
updates.</td>
</tr>
    
  
<tr><td><code>ConcurrentNamespaceSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentNamespaceSyncs is the number of namespace objects that are
allowed to sync concurrently.</td>
</tr>
    
  
</tbody>
</table>
    


## `NodeIPAMControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-NodeIPAMControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


NodeIPAMControllerConfiguration contains elements describing NodeIpamController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ServiceCIDR</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   serviceCIDR is CIDR Range for Services in cluster.</td>
</tr>
    
  
<tr><td><code>SecondaryServiceCIDR</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   secondaryServiceCIDR is CIDR Range for Services in cluster. This is used in dual stack clusters. SecondaryServiceCIDR must be of different IP family than ServiceCIDR</td>
</tr>
    
  
<tr><td><code>NodeCIDRMaskSize</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   NodeCIDRMaskSize is the mask size for node cidr in cluster.</td>
</tr>
    
  
<tr><td><code>NodeCIDRMaskSizeIPv4</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   NodeCIDRMaskSizeIPv4 is the mask size for node cidr in dual-stack cluster.</td>
</tr>
    
  
<tr><td><code>NodeCIDRMaskSizeIPv6</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   NodeCIDRMaskSizeIPv6 is the mask size for node cidr in dual-stack cluster.</td>
</tr>
    
  
</tbody>
</table>
    


## `NodeLifecycleControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-NodeLifecycleControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


NodeLifecycleControllerConfiguration contains elements describing NodeLifecycleController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>EnableTaintManager</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   If set to true enables NoExecute Taints and will evict all not-tolerating
Pod running on Nodes tainted with this kind of Taints.</td>
</tr>
    
  
<tr><td><code>NodeEvictionRate</code> <B>[Required]</B><br/>
<code>float32</code>
</td>
<td>
   nodeEvictionRate is the number of nodes per second on which pods are deleted in case of node failure when a zone is healthy</td>
</tr>
    
  
<tr><td><code>SecondaryNodeEvictionRate</code> <B>[Required]</B><br/>
<code>float32</code>
</td>
<td>
   secondaryNodeEvictionRate is the number of nodes per second on which pods are deleted in case of node failure when a zone is unhealthy</td>
</tr>
    
  
<tr><td><code>NodeStartupGracePeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   nodeStartupGracePeriod is the amount of time which we allow starting a node to
be unresponsive before marking it unhealthy.</td>
</tr>
    
  
<tr><td><code>NodeMonitorGracePeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   nodeMontiorGracePeriod is the amount of time which we allow a running node to be
unresponsive before marking it unhealthy. Must be N times more than kubelet's
nodeStatusUpdateFrequency, where N means number of retries allowed for kubelet
to post node status.</td>
</tr>
    
  
<tr><td><code>PodEvictionTimeout</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   podEvictionTimeout is the grace period for deleting pods on failed nodes.</td>
</tr>
    
  
<tr><td><code>LargeClusterSizeThreshold</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   secondaryNodeEvictionRate is implicitly overridden to 0 for clusters smaller than or equal to largeClusterSizeThreshold</td>
</tr>
    
  
<tr><td><code>UnhealthyZoneThreshold</code> <B>[Required]</B><br/>
<code>float32</code>
</td>
<td>
   Zone is treated as unhealthy in nodeEvictionRate and secondaryNodeEvictionRate when at least
unhealthyZoneThreshold (no less than 3) of Nodes in the zone are NotReady</td>
</tr>
    
  
</tbody>
</table>
    


## `PersistentVolumeBinderControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-PersistentVolumeBinderControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


PersistentVolumeBinderControllerConfiguration contains elements describing
PersistentVolumeBinderController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>PVClaimBinderSyncPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   pvClaimBinderSyncPeriod is the period for syncing persistent volumes
and persistent volume claims.</td>
</tr>
    
  
<tr><td><code>VolumeConfiguration</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-VolumeConfiguration"><code>VolumeConfiguration</code></a>
</td>
<td>
   volumeConfiguration holds configuration for volume related features.</td>
</tr>
    
  
<tr><td><code>VolumeHostCIDRDenylist</code> <B>[Required]</B><br/>
<code>[]string</code>
</td>
<td>
   VolumeHostCIDRDenylist is a list of CIDRs that should not be reachable by the
controller from plugins.</td>
</tr>
    
  
<tr><td><code>VolumeHostAllowLocalLoopback</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   VolumeHostAllowLocalLoopback indicates if local loopback hosts (127.0.0.1, etc)
should be allowed from plugins.</td>
</tr>
    
  
</tbody>
</table>
    


## `PersistentVolumeRecyclerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-PersistentVolumeRecyclerConfiguration}
    



**Appears in:**

- [VolumeConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-VolumeConfiguration)


PersistentVolumeRecyclerConfiguration contains elements describing persistent volume plugins.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>MaximumRetry</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   maximumRetry is number of retries the PV recycler will execute on failure to recycle
PV.</td>
</tr>
    
  
<tr><td><code>MinimumTimeoutNFS</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   minimumTimeoutNFS is the minimum ActiveDeadlineSeconds to use for an NFS Recycler
pod.</td>
</tr>
    
  
<tr><td><code>PodTemplateFilePathNFS</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   podTemplateFilePathNFS is the file path to a pod definition used as a template for
NFS persistent volume recycling</td>
</tr>
    
  
<tr><td><code>IncrementTimeoutNFS</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   incrementTimeoutNFS is the increment of time added per Gi to ActiveDeadlineSeconds
for an NFS scrubber pod.</td>
</tr>
    
  
<tr><td><code>PodTemplateFilePathHostPath</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   podTemplateFilePathHostPath is the file path to a pod definition used as a template for
HostPath persistent volume recycling. This is for development and testing only and
will not work in a multi-node cluster.</td>
</tr>
    
  
<tr><td><code>MinimumTimeoutHostPath</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   minimumTimeoutHostPath is the minimum ActiveDeadlineSeconds to use for a HostPath
Recycler pod.  This is for development and testing only and will not work in a multi-node
cluster.</td>
</tr>
    
  
<tr><td><code>IncrementTimeoutHostPath</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   incrementTimeoutHostPath is the increment of time added per Gi to ActiveDeadlineSeconds
for a HostPath scrubber pod.  This is for development and testing only and will not work
in a multi-node cluster.</td>
</tr>
    
  
</tbody>
</table>
    


## `PodGCControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-PodGCControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


PodGCControllerConfiguration contains elements describing PodGCController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>TerminatedPodGCThreshold</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   terminatedPodGCThreshold is the number of terminated pods that can exist
before the terminated pod garbage collector starts deleting terminated pods.
If <= 0, the terminated pod garbage collector is disabled.</td>
</tr>
    
  
</tbody>
</table>
    


## `ReplicaSetControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-ReplicaSetControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


ReplicaSetControllerConfiguration contains elements describing ReplicaSetController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentRSSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentRSSyncs is the number of replica sets that are  allowed to sync
concurrently. Larger number = more responsive replica  management, but more
CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
    


## `ReplicationControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-ReplicationControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


ReplicationControllerConfiguration contains elements describing ReplicationController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentRCSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentRCSyncs is the number of replication controllers that are
allowed to sync concurrently. Larger number = more responsive replica
management, but more CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
    


## `ResourceQuotaControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-ResourceQuotaControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


ResourceQuotaControllerConfiguration contains elements describing ResourceQuotaController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ResourceQuotaSyncPeriod</code> <B>[Required]</B><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   resourceQuotaSyncPeriod is the period for syncing quota usage status
in the system.</td>
</tr>
    
  
<tr><td><code>ConcurrentResourceQuotaSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentResourceQuotaSyncs is the number of resource quotas that are
allowed to sync concurrently. Larger number = more responsive quota
management, but more CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
    


## `SAControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-SAControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


SAControllerConfiguration contains elements describing ServiceAccountController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ServiceAccountKeyFile</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   serviceAccountKeyFile is the filename containing a PEM-encoded private RSA key
used to sign service account tokens.</td>
</tr>
    
  
<tr><td><code>ConcurrentSATokenSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentSATokenSyncs is the number of service account token syncing operations
that will be done concurrently.</td>
</tr>
    
  
<tr><td><code>RootCAFile</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   rootCAFile is the root certificate authority will be included in service
account's token secret. This must be a valid PEM-encoded CA bundle.</td>
</tr>
    
  
</tbody>
</table>
    


## `StatefulSetControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-StatefulSetControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


StatefulSetControllerConfiguration contains elements describing StatefulSetController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentStatefulSetSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentStatefulSetSyncs is the number of statefulset objects that are
allowed to sync concurrently. Larger number = more responsive statefulsets,
but more CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
    


## `TTLAfterFinishedControllerConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-TTLAfterFinishedControllerConfiguration}
    



**Appears in:**

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


TTLAfterFinishedControllerConfiguration contains elements describing TTLAfterFinishedController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentTTLSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentTTLSyncs is the number of TTL-after-finished collector workers that are
allowed to sync concurrently.</td>
</tr>
    
  
</tbody>
</table>
    


## `VolumeConfiguration`     {#kubecontrollermanager-config-k8s-io-v1alpha1-VolumeConfiguration}
    



**Appears in:**

- [PersistentVolumeBinderControllerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-PersistentVolumeBinderControllerConfiguration)


VolumeConfiguration contains &lowast;all&lowast; enumerated flags meant to configure all volume
plugins. From this config, the controller-manager binary will create many instances of
volume.VolumeConfig, each containing only the configuration needed for that plugin which
are then passed to the appropriate plugin. The ControllerManager binary is the only part
of the code which knows what plugins are supported and which flags correspond to each plugin.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>EnableHostPathProvisioning</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   enableHostPathProvisioning enables HostPath PV provisioning when running without a
cloud provider. This allows testing and development of provisioning features. HostPath
provisioning is not supported in any way, won't work in a multi-node cluster, and
should not be used for anything other than testing or development.</td>
</tr>
    
  
<tr><td><code>EnableDynamicProvisioning</code> <B>[Required]</B><br/>
<code>bool</code>
</td>
<td>
   enableDynamicProvisioning enables the provisioning of volumes when running within an environment
that supports dynamic provisioning. Defaults to true.</td>
</tr>
    
  
<tr><td><code>PersistentVolumeRecyclerConfiguration</code> <B>[Required]</B><br/>
<a href="#kubecontrollermanager-config-k8s-io-v1alpha1-PersistentVolumeRecyclerConfiguration"><code>PersistentVolumeRecyclerConfiguration</code></a>
</td>
<td>
   persistentVolumeRecyclerConfiguration holds configuration for persistent volume plugins.</td>
</tr>
    
  
<tr><td><code>FlexVolumePluginDir</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   volumePluginDir is the full path of the directory in which the flex
volume plugin should search for additional third party volume plugins</td>
</tr>
    
  
</tbody>
</table>
    
  
  
    

## `ServiceControllerConfiguration`     {#ServiceControllerConfiguration}
    



**Appears in:**

- [CloudControllerManagerConfiguration](#cloudcontrollermanager-config-k8s-io-v1alpha1-CloudControllerManagerConfiguration)

- [KubeControllerManagerConfiguration](#kubecontrollermanager-config-k8s-io-v1alpha1-KubeControllerManagerConfiguration)


ServiceControllerConfiguration contains elements describing ServiceController.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>ConcurrentServiceSyncs</code> <B>[Required]</B><br/>
<code>int32</code>
</td>
<td>
   concurrentServiceSyncs is the number of services that are
allowed to sync concurrently. Larger number = more responsive service
management, but more CPU (and network) load.</td>
</tr>
    
  
</tbody>
</table>
