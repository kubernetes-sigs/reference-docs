---
api_metadata:
  apiVersion: "apps/v1"
  import: "k8s.io/api/apps/v1"
  kind: "Deployment"
content_type: "api_reference"
description: "Deployment enables declarative updates for Pods and ReplicaSets."
title: "Deployment"
weight: 10
auto_generated: true
---

<!--
The file is auto-generated from the Go source code of the component using a generic
[generator](https://github.com/kubernetes-sigs/reference-docs/). To learn how
to generate the reference documentation, please read
[Contributing to the reference documentation](/docs/contribute/generate-ref-docs/).
To update the reference content, please follow the
[Contributing upstream](/docs/contribute/generate-ref-docs/contribute-upstream/)
guide. You can file document formatting bugs against the
[reference-docs](https://github.com/kubernetes-sigs/reference-docs/) project.
-->

`apiVersion: apps/v1`

`import "k8s.io/api/apps/v1"`

## Deployment {#Deployment}

Deployment enables declarative updates for Pods and ReplicaSets.

<hr>

<table>
  <thead><tr><th>Field</th><th>Description</th></tr></thead>
  <tbody>
    <tr>
      <td><code>apiVersion</code><br/><em>string</em><br/><em>const: <code>apps/v1</code></em></td>
      <td>APIVersion defines the versioned schema of this representation of an object.</td>
    </tr>
    <tr>
      <td><code>kind</code><br/><em>string</em><br/><em>const: <code>Deployment</code></em></td>
      <td>Kind is a string value representing the REST resource.</td>
    </tr>
    <tr>
      <td><code>metadata</code><br/><em>ObjectMeta</em></td>
      <td>Standard object's metadata.</td>
    </tr>
    <tr>
      <td><code>spec</code><br/><em>DeploymentSpec</em></td>
      <td>Specification of the desired behavior of the Deployment.</td>
    </tr>
    <tr>
      <td><code>status</code><br/><em>DeploymentStatus</em></td>
      <td>Most recently observed status of the Deployment.</td>
    </tr>
  </tbody>
</table>






