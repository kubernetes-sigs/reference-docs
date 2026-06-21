---
content_type: "api_reference"
title: "List Pods"
weight: 10
auto_generated: true
---

### `get` List Pods

#### HTTP Request

GET /api/v1/namespaces/{namespace}/pods


#### Path Parameters

<table>
  <thead><tr><th>Name</th><th>Type</th><th>Description</th></tr></thead>
  <tbody>
    <tr>
      <td><code>namespace</code></td>
      <td><em>string</em></td>
      <td>object name and auth scope, such as for teams and users</td>
    </tr>
  </tbody>
</table>


#### Query Parameters

<table>
  <thead><tr><th>Name</th><th>Type</th><th>Description</th></tr></thead>
  <tbody>
    <tr>
      <td><code>watch</code></td>
      <td><em>boolean</em></td>
      <td>Watch for changes to the described resources.</td>
    </tr>
  </tbody>
</table>


#### Body Parameters

<table>
  <thead><tr><th>Name</th><th>Type</th><th>Description</th></tr></thead>
  <tbody>
    <tr>
      <td><code>body</code></td>
      <td><em>Pod</em></td>
      <td>Pod to create.</td>
    </tr>
  </tbody>
</table>


#### Response

<table>
  <thead><tr><th>Status</th><th>Description</th><th>Response</th></tr></thead>
  <tbody>
    <tr>
      <td>200</td>
      <td>OK</td>
      <td><em>PodList</em></td>
    </tr>
    <tr>
      <td>401</td>
      <td>Unauthorized</td>
      <td>&mdash;</td>
    </tr>
  </tbody>
</table>

