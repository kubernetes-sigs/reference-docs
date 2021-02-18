---
title: kube-apiserver Encryption Configuration (v1)
content_type: tool-reference
package: apiserver.config.k8s.io/v1
auto_generated: true
---
Package v1 is the v1 version of the API.

## Resource Types 


- [EncryptionConfiguration](#apiserver-config-k8s-io-v1-EncryptionConfiguration)
  
    


## `EncryptionConfiguration`     {#apiserver-config-k8s-io-v1-EncryptionConfiguration}
    




EncryptionConfiguration stores the complete configuration for encryption providers.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>apiserver.config.k8s.io/v1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>EncryptionConfiguration</code></td></tr>
    

  
  
<tr><td><code>resources</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-ResourceConfiguration"><code>[]ResourceConfiguration</code></a>
</td>
<td>
   resources is a list containing resources, and their corresponding encryption providers.</td>
</tr>
    
  
</tbody>
</table>
    


## `AESConfiguration`     {#apiserver-config-k8s-io-v1-AESConfiguration}
    



**Appears in:**

- [ProviderConfiguration](#apiserver-config-k8s-io-v1-ProviderConfiguration)


AESConfiguration contains the API configuration for an AES transformer.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>keys</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-Key"><code>[]Key</code></a>
</td>
<td>
   keys is a list of keys to be used for creating the AES transformer.
Each key has to be 32 bytes long for AES-CBC and 16, 24 or 32 bytes for AES-GCM.</td>
</tr>
    
  
</tbody>
</table>
    


## `IdentityConfiguration`     {#apiserver-config-k8s-io-v1-IdentityConfiguration}
    



**Appears in:**

- [ProviderConfiguration](#apiserver-config-k8s-io-v1-ProviderConfiguration)


IdentityConfiguration is an empty struct to allow identity transformer in provider configuration.


    


## `KMSConfiguration`     {#apiserver-config-k8s-io-v1-KMSConfiguration}
    



**Appears in:**

- [ProviderConfiguration](#apiserver-config-k8s-io-v1-ProviderConfiguration)


KMSConfiguration contains the name, cache size and path to configuration file for a KMS based envelope transformer.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   name is the name of the KMS plugin to be used.</td>
</tr>
    
  
<tr><td><code>cachesize</code><br/>
<code>int32</code>
</td>
<td>
   cachesize is the maximum number of secrets which are cached in memory. The default value is 1000.
Set to a negative value to disable caching.</td>
</tr>
    
  
<tr><td><code>endpoint</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   endpoint is the gRPC server listening address, for example "unix:///var/run/kms-provider.sock".</td>
</tr>
    
  
<tr><td><code>timeout</code><br/>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration"><code>meta/v1.Duration</code></a>
</td>
<td>
   timeout for gRPC calls to kms-plugin (ex. 5s). The default is 3 seconds.</td>
</tr>
    
  
</tbody>
</table>
    


## `Key`     {#apiserver-config-k8s-io-v1-Key}
    



**Appears in:**

- [AESConfiguration](#apiserver-config-k8s-io-v1-AESConfiguration)

- [SecretboxConfiguration](#apiserver-config-k8s-io-v1-SecretboxConfiguration)


Key contains name and secret of the provided key for a transformer.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   name is the name of the key to be used while storing data to disk.</td>
</tr>
    
  
<tr><td><code>secret</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   secret is the actual key, encoded in base64.</td>
</tr>
    
  
</tbody>
</table>
    


## `ProviderConfiguration`     {#apiserver-config-k8s-io-v1-ProviderConfiguration}
    



**Appears in:**

- [ResourceConfiguration](#apiserver-config-k8s-io-v1-ResourceConfiguration)


ProviderConfiguration stores the provided configuration for an encryption provider.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>aesgcm</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-AESConfiguration"><code>AESConfiguration</code></a>
</td>
<td>
   aesgcm is the configuration for the AES-GCM transformer.</td>
</tr>
    
  
<tr><td><code>aescbc</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-AESConfiguration"><code>AESConfiguration</code></a>
</td>
<td>
   aescbc is the configuration for the AES-CBC transformer.</td>
</tr>
    
  
<tr><td><code>secretbox</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-SecretboxConfiguration"><code>SecretboxConfiguration</code></a>
</td>
<td>
   secretbox is the configuration for the Secretbox based transformer.</td>
</tr>
    
  
<tr><td><code>identity</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-IdentityConfiguration"><code>IdentityConfiguration</code></a>
</td>
<td>
   identity is the (empty) configuration for the identity transformer.</td>
</tr>
    
  
<tr><td><code>kms</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-KMSConfiguration"><code>KMSConfiguration</code></a>
</td>
<td>
   kms contains the name, cache size and path to configuration file for a KMS based envelope transformer.</td>
</tr>
    
  
</tbody>
</table>
    


## `ResourceConfiguration`     {#apiserver-config-k8s-io-v1-ResourceConfiguration}
    



**Appears in:**

- [EncryptionConfiguration](#apiserver-config-k8s-io-v1-EncryptionConfiguration)


ResourceConfiguration stores per resource configuration.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>resources</code> <B>[Required]</B><br/>
<code>[]string</code>
</td>
<td>
   resources is a list of kubernetes resources which have to be encrypted.</td>
</tr>
    
  
<tr><td><code>providers</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-ProviderConfiguration"><code>[]ProviderConfiguration</code></a>
</td>
<td>
   providers is a list of transformers to be used for reading and writing the resources to disk.
eg: aesgcm, aescbc, secretbox, identity.</td>
</tr>
    
  
</tbody>
</table>
    


## `SecretboxConfiguration`     {#apiserver-config-k8s-io-v1-SecretboxConfiguration}
    



**Appears in:**

- [ProviderConfiguration](#apiserver-config-k8s-io-v1-ProviderConfiguration)


SecretboxConfiguration contains the API configuration for an Secretbox transformer.

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    

  
<tr><td><code>keys</code> <B>[Required]</B><br/>
<a href="#apiserver-config-k8s-io-v1-Key"><code>[]Key</code></a>
</td>
<td>
   keys is a list of keys to be used for creating the Secretbox transformer.
Each key has to be 32 bytes long.</td>
</tr>
    
  
</tbody>
</table>
    
  
